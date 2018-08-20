package rxgo

func Let(fn func(Observable) Observable) OperatorFunc {
	return func(o Observable) Observable {
		return fn(o)
	}
}
