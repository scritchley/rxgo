package rxgo

func Take(num int) OperatorFunc {
	return func(o Observable) Observable {
		return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
			count := 0
			sub := o.Subscribe(OnNext(func(val Value) {
				if count < num {
					v <- val
				}
				if count == num-1 {
					c <- true
				}
				count++
			}).OnErr(e.Error).OnComplete(c.Complete))
			return sub.Unsubscribe
		})
	}
}

func TakeUntil(observable Observable) OperatorFunc {
	return func(o Observable) Observable {
		return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
			sub := o.Subscribe(OnNext(func(val Value) { v <- val }).OnErr(e.Error).OnComplete(c.Complete))
			observable.Subscribe(OnNext(func(val Value) { c <- true }).OnErr(func(e error) { c <- true }).OnComplete(func() { c <- true }))
			return sub.Unsubscribe
		})
	}
}
