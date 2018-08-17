package rxgo

func WithLatestFrom(other Observable) OperatorFunc {
	return func(obs Observable) Observable {
		return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
			baseObs := NewChanObserver()
			baseSub := obs.Subscribe(baseObs)
			otherObs := NewChanObserver()
			otherSub := other.Subscribe(otherObs)
			go func() {
				var a Value
				var b Value
			LOOP:
				for {
					select {
					case val := <-baseObs.ValueChan:
						a = val
						if a != nil && b != nil {
							v.Next([]Value{a, b})
						}
					case otherVal := <-otherObs.ValueChan:
						b = otherVal
						if a != nil && b != nil {
							v.Next([]Value{a, b})
						}
					case <-baseObs.CompleteChan:
						c.Complete()
						break LOOP
					case <-otherObs.CompleteChan:
						c.Complete()
						break LOOP
					case err := <-baseObs.ErrChan:
						e.Error(err)
						break LOOP
					case err := <-otherObs.ErrChan:
						e.Error(err)
						break LOOP
					}
				}
			}()
			return call(
				baseSub.Unsubscribe,
				otherSub.Unsubscribe,
			)
		})
	}
}
