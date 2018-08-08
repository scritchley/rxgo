package rxgo

func Reduce(accumulator func(acc, value Value) Value) OperatorFunc {
	return func(o Observable) Observable {
		return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
			accValue := make(chan Value, 1)
			return o.Subscribe(
				OnNext(func(val Value) {
					select {
					case acc := <-accValue:
						accValue <- accumulator(acc, val)
					default:
						accValue <- val
					}
				}).OnComplete(func() {
					v <- <-accValue
					c.Complete()
				}),
			).Unsubscribe
		})
	}
}
