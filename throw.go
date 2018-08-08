package rxgo

func Throw(err error) Observable {
	return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
		e.Error(err)
		return nil
	})
}
