package rxgo

func Multicast(fn func() Subject) OperatorFunc {
	return func(o Observable) Observable {
		subject := fn()
		onSubscribe := make(chan bool, 1)
		go func() {
			<-onSubscribe
			o.Subscribe(subject)
		}()
		return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
			select {
			case onSubscribe <- true:
			default:
			}
			return subject.Subscribe(OnNext(v.Next).OnErr(e.Error).OnComplete(c.Complete)).Unsubscribe
		})
	}
}

func Share() OperatorFunc {
	return Multicast(NewSubject)
}

func ShareReplay(buffer int) OperatorFunc {
	return Multicast(func() Subject {
		return NewReplaySubject(buffer)
	})
}
