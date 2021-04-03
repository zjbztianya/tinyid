package fanout

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrFull = errors.New("fanout buffer full")
)

type task struct {
	fn  func(ctx context.Context)
	ctx context.Context
}

type Fanout struct {
	ch         chan task
	buffer     int
	numsWorker int
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewFanout(buffer int, numsWorker int) *Fanout {
	ctx, cancel := context.WithCancel(context.Background())
	fanout := &Fanout{
		buffer:     buffer,
		numsWorker: numsWorker,
		ctx:        ctx,
		cancel:     cancel,
		ch:         make(chan task, buffer)}

	fanout.wg.Add(numsWorker)
	for i := 0; i < numsWorker; i++ {
		go fanout.proc()
	}

	return fanout
}

func (f *Fanout) proc() {
	defer f.wg.Done()
	for {
		select {
		case t := <-f.ch:
			t.fn(t.ctx)
		case <-f.ctx.Done():
			return
		}
	}
}

func (f *Fanout) Do(ctx context.Context, fn func(context.Context)) error {
	if fn == nil || f.ctx.Err() != nil {
		return f.ctx.Err()
	}
	select {
	case f.ch <- task{
		fn:  fn,
		ctx: ctx,
	}:
	default:
		return ErrFull
	}
	return nil
}

func (f *Fanout) Close() error {
	if err := f.ctx.Err(); err != nil {
		return err
	}
	f.cancel()
	f.wg.Wait()
	return nil
}
