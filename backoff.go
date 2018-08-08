package rxgo

import (
	"time"
)

func Backoff(d time.Duration) OperatorFunc {
	return ConcatMap(func(v Value) Observable {
		return Of(v).Pipe(Delay(d), Do(func(Value) { d = d * 2 }))
	})
}

func BackoffWithMax(d, max time.Duration) OperatorFunc {
	return ConcatMap(func(v Value) Observable {
		return Of(v).Pipe(Delay(d), Do(func(Value) {
			d = d * 2
			if d > max {
				d = max
			}
		}))
	})
}
