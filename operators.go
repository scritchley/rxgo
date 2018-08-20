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

// ToSlice is an alias of ToArray with a more Go-friendly name.
func ToSlice() OperatorFunc {
	return ToArray()
}

func ToArray() OperatorFunc {
	return func(o Observable) Observable {
		return o.Pipe(
			Reduce(func(acc, val Value) Value {
				return append(acc.([]Value), val)
			}, []Value(nil)),
		)
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

func Last() OperatorFunc {
	return func(o Observable) Observable {
		return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
			var lastValue Value
			return o.Subscribe(OnNext(func(val Value) {
				lastValue = val
			}).OnErr(func(err error) {
				v.Next(lastValue)
				e.Error(err)
			}).OnComplete(func() {
				v.Next(lastValue)
				c.Complete()
			})).Unsubscribe
		})
	}
}
