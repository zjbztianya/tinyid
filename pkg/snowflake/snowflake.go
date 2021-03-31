package snowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	//Milliseconds since the epoch,We use Twitter snowflake default epoch.
	TwEpoch           = 1288834974657
	WorkerBits        = 10
	MaxWorkerID       = (1 << WorkerBits) - 1
	SequenceBits      = 12
	TimestampLeftShit = SequenceBits + WorkerBits
	SequenceMask      = (1 << SequenceBits) - 1
	WaitMills         = 100 * time.Microsecond
)

// SnowFlake is Twitter ID generator
// reference:https://blog.twitter.com/engineering/en_us/a/2010/announcing-snowflake.html
type SnowFlake struct {
	mu            sync.Mutex
	lastTimestamp int64
	workerID      int64
	sequence      int64
}

func NewSnowFlake(workerID int64) *SnowFlake {
	if workerID < 0 || workerID > MaxWorkerID {
		return nil
	}
	return &SnowFlake{workerID: workerID}
}

func (s *SnowFlake) GenID() (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixNano() / 1e6
	if now == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & SequenceMask
		if s.sequence == 0 {
			now = s.waitNextMillis(now)
		}
	} else {
		s.sequence = 0
	}

	// avoid clock backwards
	if now < s.lastTimestamp {
		return 0, errors.New("inner time error")
	}

	s.lastTimestamp = now
	return s.generate(), nil
}

func (s *SnowFlake) generate() int64 {
	return ((s.lastTimestamp - TwEpoch) << TimestampLeftShit) | (s.workerID << SequenceBits) | s.sequence
}

func (s *SnowFlake) waitNextMillis(t int64) int64 {
	for t == s.lastTimestamp {
		time.Sleep(WaitMills)
		t = time.Now().UnixNano() / 1e6
	}
	return t
}