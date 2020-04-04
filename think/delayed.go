package think

import (
	"context"
	"sync"
	"time"

	"github.com/francoispqt/onelog"
)

type DelayedCall struct {
	wg     *sync.WaitGroup
	ctx    *context.Context
	cancel context.CancelFunc
	action func()
	logger *onelog.Logger
}

func (call *DelayedCall) Wait() {
	call.logger.Debug("waiting for delayed action")
	call.wg.Wait()
}

// Do blocks until `Run` is called or time is out and the runs `action`
func (call *DelayedCall) Do() {
	<-(*call.ctx).Done()
	call.logger.Debug("running delayed action")
	call.action()
	call.wg.Done()
	call.logger.Debug("delayed action is done")
}

// Run immediately runs `action` and waits until it is done
func (call *DelayedCall) Run() {
	call.logger.Debug("running delayed action")
	call.cancel()
	call.Wait()
}

type DelayedCalls struct {
	queue  chan *DelayedCall
	logger *onelog.Logger
}

func (calls *DelayedCalls) Add(duration time.Duration, jobID int, action func()) *DelayedCall {
	logger := calls.logger.With(func(e onelog.Entry) {
		e.Int("job", jobID)
	})
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	wg := sync.WaitGroup{}
	wg.Add(1)
	call := DelayedCall{
		wg:     &wg,
		ctx:    &ctx,
		cancel: cancel,
		action: action,
		logger: logger,
	}
	calls.queue <- &call
	return &call
}

// Run runs all delayed actions at once
func (calls *DelayedCalls) Run() {
	for {
		select {
		case call := <-calls.queue:
			call.Run()
		default:
			return
		}
	}
}
