package rxgo

import (
	"time"
)

func Delay(d time.Duration) OperatorFunc {
	return func(o Observable) Observable {
		return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
			return o.Subscribe(OnNext(func(val Value) {
				time.Sleep(d)
				v.Next(val)
			}).OnErr(e.Error).OnComplete(c.Complete)).Unsubscribe
		})
	}
}

func Debounce(d time.Duration) OperatorFunc {
	return func(o Observable) Observable {
		return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
			debounce := make(chan Value)
			done := make(chan bool)
			go func() {
			LOOP:
				for {
					time.Sleep(d)
					v.Next(<-debounce)
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
				c.Complete()
			})).Unsubscribe
		})
	}
}

func Buffer(num int) OperatorFunc {
	return BufferWithCount(num, num)
}

func BufferWithCount(num, count int) OperatorFunc {
	if count < 1 {
		panic("count must be greater than 0")
	}
	return func(o Observable) Observable {
		return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
			var buffer []Value
			return o.Subscribe(OnNext(func(val Value) {
				buffer = append(buffer, val)
				if len(buffer) >= num && len(buffer)%count == 0 {
					v.Next(buffer[len(buffer)-num:])
					buffer = buffer[len(buffer)-num:]
				}
			}).OnErr(e.Error).OnComplete(c.Complete)).Unsubscribe
		})
	}
}

func ToArray() OperatorFunc {
	return func(o Observable) Observable {
		return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
			var values []Value
			return o.Subscribe(OnNext(func(val Value) {
				values = append(values, val)
			}).OnErr(e.Error).OnComplete(func() {
				v.Next(values)
				c.Complete()
			})).Unsubscribe
		})
	}
}

func PairWise() OperatorFunc {
	return BufferWithCount(2, 1)
}

func First(predicate func(v Value) bool) OperatorFunc {
	if predicate == nil {
		return Take(1)
	}
	return Filter(predicate).Pipe(Take(1))
}
