package rxgo

import "time"

func Timeout(timeout time.Duration) OperatorFunc {
	return func(o Observable) Observable {
		first := NewSubject()
		cancelled := NewSubject()
		return Merge(
			o.Pipe(
				Do(func(v Value) {
					first.Complete()
				}),
				TakeUntil(cancelled),
			),
			Of(1).Pipe(
				Delay(timeout),
				Do(func(v Value) {
					cancelled.Complete()
				}),
				TakeUntil(first),
				IgnoreElements(),
			),
		)
	}
}
