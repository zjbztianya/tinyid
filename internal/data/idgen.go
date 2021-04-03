package data

import (
	"context"
	"database/sql"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"golang.org/x/sync/singleflight"
	"sync"
	"sync/atomic"
	"time"
	"tinyid/internal/biz"
	"tinyid/pkg/fanout"
)

const (
	_getAllTinyidAlloc       = "SELECT biz_tag, max_id, step, update_time FROM tinyid_alloc"
	_getTinyidAlloc          = "SELECT biz_tag, max_id, step FROM tinyid_alloc WHERE biz_tag = ?"
	_updateMaxID             = "UPDATE tinyid_alloc SET max_id = max_id + step WHERE biz_tag = ?"
	_updateMaxIDByCustomStep = "UPDATE tinyid_alloc SET max_id = max_id + ? WHERE biz_tag = ?"
	_getAllTags              = "SELECT biz_tag FROM tinyid_alloc"
	segmentDuration          = 15 * time.Minute
	maxStep                  = 1000000
)

type IDAlloc struct {
	BizTag    string
	MaxID     int64
	Step      int
	UpdatedAt time.Time
}

type idgenRepo struct {
	data *Data
	log  *log.Helper
	// segmentID local cache
	segmentCache sync.Map
	flight       *singleflight.Group
	fanout       *fanout.Fanout
}

// NewIdgenRepo .
func NewIdgenRepo(data *Data, logger log.Logger) biz.IdgenRepo {
	repo := &idgenRepo{
		data:   data,
		log:    log.NewHelper("data/idgen", logger),
		flight: &singleflight.Group{},
		fanout: fanout.NewFanout(8, 2),
	}

	repo.updateSegmentCache()
	go repo.loadSegmentCacheJob()
	return repo
}

func (r *idgenRepo) WorkerID() int64 {
	return r.data.workerID
}

func (r *idgenRepo) loadSegmentCacheJob() {
	tick := time.NewTicker(60 * time.Second)
	for range tick.C {
		r.updateSegmentCache()
	}
}

func (r *idgenRepo) updateSegmentCache() {
	var (
		dbTags []string
		err    error
	)

	dbTagMap := make(map[string]struct{})
	for {
		dbTags, err = r.tags(context.Background())
		if err != nil {
			r.log.Debug(err)
			time.Sleep(time.Second)
			continue
		}

		//Add new biz tag cache
		for _, tag := range dbTags {
			dbTagMap[tag] = struct{}{}
			if _, ok := r.segmentCache.Load(tag); !ok {
				r.segmentCache.Store(tag, NewSegmentBuffer())
			}
		}

		//Delete invalid tag cache
		r.segmentCache.Range(func(key, value interface{}) bool {
			if _, ok := dbTagMap[key.(string)]; !ok {
				r.segmentCache.Delete(key)
			}
			return true
		})
		return
	}
}

func (r *idgenRepo) updateSegment(tag string, idx int) (err error) {
	var (
		idAlloc  *IDAlloc
		nextStep int
	)
	b, ok := r.segmentCache.Load(tag)
	if !ok {
		return err
	}

	buffer := b.(*SegmentBuffer)
	db := r.data.db
	ctx, cancel := context.WithTimeout(context.Background(), db.conf.TranTimeout.AsDuration())
	defer cancel()

	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	rollback := func() {
		if err1 := tx.Rollback(); err1 != nil {
			r.log.Errorf("tx.Rollback() error(%v)", err1)
		}
	}

	if atomic.LoadInt32(&buffer.initOK) == 0 {
		err = r.txSetMaxID(ctx, tx, tag)
		if err != nil {
			rollback()
			return
		}
		idAlloc, err = r.txAlloc(ctx, tx, tag)
		if err != nil {
			rollback()
			return
		}
		if err = tx.Commit(); err != nil {
			return errors.Wrap(err, "tx.Commit() failed")
		}

		nextStep = idAlloc.Step
	} else {
		buffer.mu.RLock()
		duration := time.Duration(time.Now().UnixNano() - buffer.lastTime)
		nextStep = buffer.step
		if duration < segmentDuration {
			if nextStep*2 < maxStep {
				nextStep *= 2
			}
		} else if duration > 2*segmentDuration {
			minStep := buffer.minStep
			if nextStep/2 >= minStep {
				nextStep /= 2
			}
		}
		buffer.mu.RUnlock()

		err = r.txSetMaxIDByCustomStep(ctx, tx, tag, int(nextStep))
		if err != nil {
			rollback()
			return
		}
		idAlloc, err = r.txAlloc(ctx, tx, tag)
		if err != nil {
			rollback()
			return
		}
		if err = tx.Commit(); err != nil {
			return errors.Wrap(err, "tx.Commit() failed")
		}
	}

	buffer.mu.Lock()
	defer buffer.mu.Unlock()

	buffer.step = nextStep
	buffer.minStep = idAlloc.Step
	buffer.lastTime = time.Now().UnixNano()
	if buffer.curIdx != idx {
		buffer.nextReady = 1
	}

	segment := buffer.segments[idx]
	segment.value = idAlloc.MaxID - int64(buffer.step)
	segment.maxID = idAlloc.MaxID
	segment.step = buffer.step
	r.log.Debugf("update segment success,%v", segment)
	return nil
}

func (r *idgenRepo) segmentID(ctx context.Context, buffer *SegmentBuffer, tag string) (int64, error) {
	buffer.mu.RLock()

	segment := buffer.current()
	r.log.Debugf("segment current idx:%d", buffer.curIdx)
	if buffer.nextReady == 0 &&
		float64(segment.idle()) < 0.9*float64(segment.step) &&
		atomic.CompareAndSwapInt32(&buffer.loading, 0, 1) {
		r.log.Debug("async load segment")
		// 异步加载段号
		r.fanout.Do(ctx, func(ctx context.Context) {
			_ = r.updateSegment(tag, (buffer.curIdx+1)%2)
			atomic.CompareAndSwapInt32(&buffer.loading, 1, 0)
		})
	}

	value := atomic.AddInt64(&segment.value, int64(1))
	r.log.Debugf("add segment value:%d maxID:%d", value, segment.maxID)
	if value <= segment.maxID {
		buffer.mu.RUnlock()
		return value, nil
	}

	buffer.mu.RUnlock()

	buffer.mu.Lock()
	defer buffer.mu.Unlock()
	if buffer.nextReady == 0 {
		return 0, errors.Errorf("both two segment buffer not ready:%s", tag)
	}
	segment.value++
	buffer.curIdx = (buffer.curIdx + 1) % 2
	buffer.nextReady = 0
	return segment.value, nil
}

func (r *idgenRepo) SegmentID(ctx context.Context, tag string) (int64, error) {
	b, ok := r.segmentCache.Load(tag)
	if !ok {
		return 0, errors.Errorf("segment cache miss,tag:%s", tag)
	}

	buffer := b.(*SegmentBuffer)
	if atomic.LoadInt32(&buffer.initOK) == 0 {
		r.log.Debugf("start init buffer")
		_, err, _ := r.flight.Do(tag, func() (interface{}, error) {
			err := r.updateSegment(tag, 0)
			return nil, err
		})
		if err != nil {
			r.log.Errorf("init buffer error:%v", err)
			return 0, err
		}
		r.log.Debugf("init buffer ok!")
		atomic.CompareAndSwapInt32(&buffer.initOK, 0, 1)
	}
	return r.segmentID(ctx, buffer, tag)
}

func (r *idgenRepo) tags(ctx context.Context) ([]string, error) {
	rows, err := r.data.db.db.QueryContext(ctx, _getAllTags)
	if err != nil {
		return nil, errors.Wrap(err, "tags:db.Query")
	}

	defer rows.Close()
	var res []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, errors.Wrap(err, "tags:row.Scan row")
		}
		res = append(res, tag)
	}
	return res, nil
}

func (r *idgenRepo) txAlloc(ctx context.Context, tx *sql.Tx, tag string) (alloc *IDAlloc, err error) {
	alloc = new(IDAlloc)
	row := tx.QueryRowContext(ctx, _getTinyidAlloc, tag)
	if err = row.Scan(&alloc.BizTag, &alloc.MaxID, &alloc.Step); err != nil {
		return nil, errors.Wrapf(err, "row.Scan error,tag:%s", tag)
	}
	return
}

func (r *idgenRepo) txSetMaxID(ctx context.Context, tx *sql.Tx, tag string) error {
	_, err := tx.ExecContext(ctx, _updateMaxID, tag)
	if err != nil {
		return errors.Wrap(err, "tx.ExecContext error")
	}
	return err
}

func (r *idgenRepo) txSetMaxIDByCustomStep(ctx context.Context, tx *sql.Tx, tag string, step int) error {
	_, err := tx.ExecContext(ctx, _updateMaxIDByCustomStep, step, tag)
	if err != nil {
		return errors.Wrap(err, "tx.ExecContext error")
	}
	return err
}
