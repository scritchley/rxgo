package rxgo

func Catch(catch func(err error) (Value, error)) OperatorFunc {
	return func(o Observable) Observable {
		return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
			return o.Subscribe(OnNext(v.Next).OnErr(func(err error) {
				val, nextErr := catch(err)
				if nextErr != nil {
					e.Error(nextErr)
				} else if val != nil {
					v.Next(val)
					c.Complete()
				} else {
					c.Complete()
				}
			}).OnComplete(c.Complete)).Unsubscribe
		})
	}
}
