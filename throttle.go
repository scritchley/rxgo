package rxgo

import "time"

func Throttle(durationSelector func(Value)) OperatorFunc {
	return func(o Observable) Observable {
		return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
			debounce := make(chan Value)
			done := make(chan bool)
			go func() {
			LOOP:
				for {
					val := <-debounce
					durationSelector(val)
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
	return Throttle(func(Value) {
		time.Sleep(d)
	})
}
