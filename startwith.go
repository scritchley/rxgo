package rxgo

func StartWith(value Value) OperatorFunc {
	return func(o Observable) Observable {
		return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
			firstValue := make(chan Value, 1)
			firstValue <- value
			return o.Subscribe(OnNext(func(val Value) {
				select {
				case firstVal := <-firstValue:
					v.Next(firstVal)
				default:
				}
				v.Next(val)
			}).OnErr(e.Error).OnComplete(c.Complete)).Unsubscribe
		})
	}
}

func EndWith(value Value) OperatorFunc {
	return func(o Observable) Observable {
		return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
			return o.Subscribe(OnNext(v.Next).OnErr(e.Error).OnComplete(func() {
				v.Next(value)
				c.Complete()
			})).Unsubscribe
		})
	}
}
