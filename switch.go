package rxgo

func Switch(observables ...Observable) Observable {
	return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {

		return nil
	})
}
