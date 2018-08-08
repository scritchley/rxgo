package rxgo

type observer struct {
	ValueChan
	ErrChan
	CompleteChan
}

type OnNext func(v Value)

type OnErr func(err error)

type OnComplete func()

func (o OnComplete) Next(Value) {
}

func (o OnComplete) Err(error) {
}

func (o OnComplete) Complete() {
	o()
}

func (o OnNext) Next(v Value) {
	o(v)
}

func (o OnNext) Err(err error) {}

func (o OnNext) Complete() {}

func (o OnNext) OnErr(oe OnErr) ObserverFuncs {
	return ObserverFuncs{
		next: o,
		err:  oe,
	}
}

func (o OnNext) OnComplete(oc OnComplete) ObserverFuncs {
	return ObserverFuncs{
		next:     o,
		complete: oc,
	}
}

func (o ObserverFuncs) Complete() {
	if o.complete != nil {
		o.complete()
	}
}

func (o ObserverFuncs) OnComplete(oc OnComplete) ObserverFuncs {
	o.complete = oc
	return o
}

func (o ObserverFuncs) Err(err error) {
	if o.err != nil {
		o.err(err)
	}
}

func (o ObserverFuncs) OnErr(oe OnErr) ObserverFuncs {
	o.err = oe
	return o
}

func (o ObserverFuncs) Next(v Value) {
	if o.next != nil {
		o.next(v)
	}
}

func (o ObserverFuncs) OnNext(on OnNext) ObserverFuncs {
	o.next = on
	return o
}

type ObserverFuncs struct {
	next     OnNext
	err      OnErr
	complete OnComplete
}

type chanObserver struct {
	ValueChan
	ErrChan
	CompleteChan
}

func NewChanObserver() chanObserver {
	return chanObserver{
		make(ValueChan),
		make(ErrChan, 1),
		make(CompleteChan, 1),
	}
}

func (c chanObserver) Next(v Value) {
	c.ValueChan <- v
}

func (c chanObserver) Err(err error) {
	c.ErrChan <- err
}

func (c chanObserver) Complete() {
	c.CompleteChan <- true
}

func (co chanObserver) From() Observable {
	return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
		go func() {
		LOOP:
			for {
				select {
				case val := <-co.ValueChan:
					v <- val
				case err := <-co.ErrChan:
					e <- err
					break LOOP
				case <-co.CompleteChan:
					c <- true
					break LOOP
				}
			}
		}()
		return func() {}
	})
}
