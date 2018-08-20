package rxgo

func OnSubscribe(fn func()) OperatorFunc {
	return func(o Observable) Observable {
		return Of(nil).Pipe(
			Do(func(v Value) { fn() }),
			ConcatMapTo(o),
		)
	}
}
