package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"tinyid/internal/biz"
)

const (
	_getAllTinyidAllocs      = "SELECT biz_tag, max_id, step, update_time FROM tinyid_alloc"
	_getTinyidAlloc          = "SELECT biz_tag, max_id, step FROM tinyid_alloc WHERE biz_tag = ?"
	_updateMaxId             = "UPDATE tinyid_alloc SET max_id = max_id + step WHERE biz_tag = ?"
	_updateMaxIdByCustomStep = "UPDATE tinyid_alloc SET max_id = max_id + ? WHERE biz_tag = ?"
	_getAllTags              = "SELECT biz_tag FROM tinyid_alloc"
)

type idgenRepo struct {
	data *Data
	log  *log.Helper
}

// NewIdgenRepo .
func NewIdgenRepo(data *Data, logger log.Logger) biz.IdgenRepo {
	return &idgenRepo{
		data: data,
		log:  log.NewHelper("data/idgen", logger),
	}
}

func (r *idgenRepo) CreateIdgen(g *biz.Idgen) error {
	return nil
}

func (r *idgenRepo) UpdateIdgen(g *biz.Idgen) error {
	return nil
}
