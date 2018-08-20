package rxgo

func IgnoreElements() OperatorFunc {
	return func(o Observable) Observable {
		return o.Pipe(
			Filter(func(Value) bool { return false }),
		)
	}
}
