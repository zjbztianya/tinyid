package data

import "sync"

type Segment struct {
	maxID int64
	step  int
	value int64
}

func (s *Segment) idle() int64 {
	return s.maxID - s.value
}

type SegmentBuffer struct {
	mu        sync.RWMutex
	segments  []*Segment
	curIdx    int
	nextReady int32
	initOK    int32
	loading   int32
	step      int
	minStep   int
	lastTime  int64
}

func NewSegmentBuffer() *SegmentBuffer {
	segments := make([]*Segment, 2)
	for i := 0; i < len(segments); i++ {
		segments[i] = new(Segment)
	}
	buffer := &SegmentBuffer{
		segments: segments,
	}
	return buffer
}

func (s *SegmentBuffer) current() *Segment {
	return s.segments[s.curIdx]
}
