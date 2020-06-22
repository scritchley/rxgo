package rxgo

import (
	"time"
)

func WindowTime(d time.Duration, apply ...OperatorFunc) OperatorFunc {
	return Window(Interval(d), apply...)
}

func Window(boundaries Observable, apply ...OperatorFunc) OperatorFunc {
	return func(o Observable) Observable {
		var window []Value
		var end bool
		o.Pipe(
			Do(func(v Value) {
				window = append(window, v)
			}),
			Finalize(func() {
				end = true
			}),
		).Subscribe(noOpObserver{})
		return o.Pipe(
			Let(func(o Observable) Observable {
				return boundaries.Pipe(
					TakeWhile(func(Value) bool {
						return !end
					}, true),
					ExhaustMap(func(v Value) Observable {
						oldWindow := window
						window = make([]Value, 0)
						return Of(oldWindow...).Pipe(apply...).Pipe(ToArray())
					}),
				)
			}),
		)
	}
}
