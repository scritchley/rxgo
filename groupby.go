package rxgo

import (
	"sync"
)

func GroupBy(fn func(v Value) Value, apply ...OperatorFunc) OperatorFunc {
	return func(o Observable) Observable {
		return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
			var wg sync.WaitGroup
			groups := make(map[Value]Subject)
			return o.Subscribe(
				OnNext(
					func(val Value) {
						key := fn(val)
						if group, ok := groups[key]; ok {
							group.Next(val)
						} else {
							wg.Add(1)
							groups[key] = NewSubject()
							groups[key].
								Pipe(apply...).
								Subscribe(
									OnNext(v.Next).
										OnErr(e.Error).
										OnComplete(wg.Done),
								)
							groups[key].Next(val)
						}
					},
				).
					OnErr(e.Error).
					OnComplete(func() {
						for _, group := range groups {
							group.Complete()
						}
						wg.Wait()
						c.Complete()
					}),
			).Unsubscribe
		})
	}
}
