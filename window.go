package rxgo

import (
	"sync"
	"time"
)

func WindowTime(d time.Duration) OperatorFunc {
	return Window(Interval(d))
}

func Window(boundaries Observable) OperatorFunc {
	return func(o Observable) Observable {
		return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
			var buffer []Value
			var mtx sync.Mutex
			done := make(chan bool)
			sub := NewSubscription(nil)
			sub.Add(boundaries.Subscribe(OnNext(func(_ Value) {
				mtx.Lock()
				defer mtx.Unlock()
				v <- buffer
				buffer = nil
				select {
				case <-done:
					c.Complete()
				default:
				}
			}).OnComplete(func() {})).Unsubscribe)
			sub.Add(o.Subscribe(OnNext(func(val Value) {
				mtx.Lock()
				defer mtx.Unlock()
				buffer = append(buffer, val)
			}).OnErr(e.Error).OnComplete(func() {
				done <- true
			})).Unsubscribe)
			return sub.Unsubscribe
		})
	}
}
