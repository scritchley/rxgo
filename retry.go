package rxgo

func Retry(maxRetries int) OperatorFunc {
	return func(o Observable) Observable {
		retryCount := 0
		return Range(0, maxRetries).Pipe(
			ConcatMapTo(
				o.Pipe(
					Catch(func(err error) (Value, error) {
						if retryCount == maxRetries {
							return nil, err
						}
						retryCount++
						return nil, nil
					}),
				),
			),
		)
	}
}
