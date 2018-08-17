package rxgo

import "time"

func Throttle(d time.Duration, scheduler Scheduler) OperatorFunc {
	return func(o Observable) Observable {
		return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
			debounce := make(chan Value)
			done := make(chan bool)
			go func() {
			LOOP:
				for {
					val := <-debounce
					scheduler.Schedule(func() {}, d)
					v <- val
					select {
					case <-done:
						break LOOP
					default:
					}
				}
				c.Complete()
			}()
			return o.Subscribe(OnNext(func(val Value) {
				select {
				case debounce <- val:
				default:
				}
			}).OnErr(e.Error).OnComplete(func() {
				done <- true
			})).Unsubscribe
		})
	}
}

func ThrottleTime(d time.Duration) OperatorFunc {
	return Throttle(d, DefaultAsyncScheduler)
}
