package rxgo

func Retry(retries int) OperatorFunc {
	return func(o Observable) Observable {
		return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
			go func() {
				for i := 0; i <= retries; i++ {
					o.Subscribe(
						OnNext(v.Next).
							OnErr(func(err error) {
								if i == retries {
									e.Error(err)
								}
							}).
							OnComplete(c.Complete),
					).Wait()
				}
			}()
			return nil
		})
	}
}
