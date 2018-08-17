package rxgo

import (
	"sync"
)

// Map applies the provided mapping func to each value emitted by the
// source observable and returns a new observable of the resulting values.
func Map(fn func(v Value) Value) OperatorFunc {
	return func(o Observable) Observable {
		return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
			return o.Subscribe(
				OnNext(func(val Value) {
					v <- fn(val)
				}).OnErr(func(err error) {
					e <- err
				}).OnComplete(func() {
					c <- true
				}),
			).Unsubscribe
		})
	}
}

// MapTo returns an observable that emits the provided value for each value
// emitted by the source observable.
func MapTo(value Value) OperatorFunc {
	return func(o Observable) Observable {
		return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
			return o.Subscribe(
				OnNext(func(val Value) {
					v <- value
				}).OnErr(func(err error) {
					e <- err
				}).OnComplete(func() {
					c <- true
				}),
			).Unsubscribe
		})
	}
}

// MergeMap applies the provided func to each value emitted by the source observable
// and subscribes to the resulting observable. It returns a new observable that merges
// the emitted values from all the inner observables and that completes onces all inner
// observables have completed.
func MergeMap(fn func(v Value) Observable) OperatorFunc {
	return func(o Observable) Observable {
		return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
			psub := NewSubscription(nil)
			var wg sync.WaitGroup
			sub := o.Subscribe(
				OnNext(
					func(val Value) {
						wg.Add(1)
						csub := fn(val).
							Subscribe(
								OnNext(v.Next).
									OnErr(e.Error).
									OnComplete(wg.Done),
							)
						psub.Add(csub.Unsubscribe)
					},
				).
					OnErr(e.Error).
					OnComplete(func() {
						wg.Wait()
						c.Complete()
					}),
			)
			psub.Add(sub.Unsubscribe)
			return psub.Unsubscribe
		})
	}
}

// FlatMap is an alias of MergeMap.
func FlatMap(fn func(v Value) Observable) OperatorFunc {
	return MergeMap(fn)
}

func ConcatMap(fn func(v Value) Observable) OperatorFunc {
	return func(o Observable) Observable {
		return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
			return o.Subscribe(
				OnNext(func(val Value) {
					fn(val).Subscribe(OnNext(v.Next).OnErr(e.Error)).Wait()
				}).OnErr(e.Error).OnComplete(c.Complete),
			).Unsubscribe
		})
	}
}

func ConcatMapTo(o Observable) OperatorFunc {
	return ConcatMap(func(v Value) Observable { return o })
}

func SwitchMap(fn func(v Value) Observable) OperatorFunc {
	return func(o Observable) Observable {
		return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
			var wg sync.WaitGroup
			next := NewSubject()
			return o.Subscribe(
				OnNext(func(val Value) {
					wg.Add(1)
					next.Next(nil)
					fn(val).Pipe(TakeUntil(next)).Subscribe(OnNext(v.Next).OnErr(e.Error).OnComplete(wg.Done))
				}).OnErr(e.Error).OnComplete(call(wg.Wait, c.Complete)),
			).Unsubscribe
		})
	}
}

func SwitchMapTo(o Observable) OperatorFunc {
	return SwitchMap(func(v Value) Observable { return o })
}
