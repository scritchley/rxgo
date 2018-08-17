package rxgo

type Connector chan bool

func (c Connector) Connect() {
	c <- true
}

func Multicast(connector Connector, fn func() Subject) OperatorFunc {
	return func(o Observable) Observable {
		subject := fn()
		go func() {
			<-connector
			o.Subscribe(subject)
		}()
		return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
			return subject.Subscribe(OnNext(v.Next).OnErr(e.Error).OnComplete(c.Complete)).Unsubscribe
		})
	}
}
