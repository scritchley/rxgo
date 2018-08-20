package rxgo

func Empty() Observable {
	return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
		c.Complete()
		return nil
	})
}
