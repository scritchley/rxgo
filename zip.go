package rxgo

func Zip(observables ...Observable) Observable {
	return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
		var observers []chanObserver
		for _, observable := range observables {
			co := NewChanObserver()
			observers = append(observers, co)
			observable.Subscribe(co)
		}
		go func() {
		LOOP:
			for {
				var values []Value
				for _, observer := range observers {
					select {
					case val := <-observer.ValueChan:
						values = append(values, val)
					case <-observer.CompleteChan:
						c.Complete()
						break LOOP
					case err := <-observer.ErrChan:
						e.Error(err)
						break LOOP
					}
				}
				v.Next(values)
			}
		}()
		return nil
	})
}
