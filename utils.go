package rxgo

func Do(do func(Value)) OperatorFunc {
	return func(o Observable) Observable {
		return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
			return o.Subscribe(OnNext(func(val Value) {
				do(val)
				v.Next(val)
			}).OnErr(e.Error).OnComplete(c.Complete)).Unsubscribe
		})
	}
}

func Finalize(finalize func()) OperatorFunc {
	return func(o Observable) Observable {
		return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
			return o.Subscribe(OnNext(v.Next).OnErr(e.Error).OnComplete(func() {
				finalize()
				c.Complete()
			})).Unsubscribe
		})
	}
}
