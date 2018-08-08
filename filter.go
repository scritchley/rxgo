package rxgo

// Filter applies the predicate func to values.
func Filter(predicate func(v Value) bool) OperatorFunc {
	return func(o Observable) Observable {
		return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
			return o.Subscribe(OnNext(func(val Value) {
				if predicate(val) {
					v <- val
				}
			}).OnErr(e.Error).OnComplete(c.Complete)).Unsubscribe
		})
	}
}
