package rxgo

type BaseObservable struct {
	Subscriber
}

func Create(subscriber Subscriber) BaseObservable {
	return BaseObservable{
		subscriber,
	}
}

func (c BaseObservable) Subscribe(obs Observer) Subscription {
	valueCh := make(ValueChan)
	errCh := make(ErrChan, 1)
	completeCh := make(CompleteChan)
	td := c.Subscriber(valueCh, errCh, completeCh)
	sub := NewSubscription(td)
	go func() {
	LOOP:
		for {
			select {
			case v := <-valueCh:
				obs.Next(v)
			case e := <-errCh:
				obs.Err(e)
				break LOOP
			case <-completeCh:
				obs.Complete()
				break LOOP
			}
		}
		sub.Unsubscribe()
	}()
	return sub
}

func (c BaseObservable) Pipe(ops ...OperatorFunc) Observable {
	return Pipe(ops...)(c)
}

type TeardownFunc func()

type Subscriber func(ValueChan, ErrChan, CompleteChan) TeardownFunc

type cancellableObserver struct {
	Observer
	cancel chan bool
}

func (c cancellableObserver) Cancel() {
	c.cancel <- true
}

type Observable interface {
	Subscribable
	Pipeable
}

type Subscribable interface {
	Subscribe(Observer) Subscription
}

type Observer interface {
	NextObserver
	ErrObserver
	CompletionObserver
}

type Value interface{}

type ValueChan chan Value

func (v ValueChan) Next(val Value) {
	v <- val
}

type ErrChan chan error

func (e ErrChan) Error(err error) {
	e <- err
}

type CompleteChan chan bool

func (c CompleteChan) Complete() {
	close(c)
}

type NextObserver interface {
	Next(Value)
}

type ErrObserver interface {
	Err(error)
}

type CompletionObserver interface {
	Complete()
}

type Pipeable interface {
	Pipe(...OperatorFunc) Observable
}

type OperatorFunc func(o Observable) Observable

func (o OperatorFunc) Pipe(fns ...OperatorFunc) OperatorFunc {
	return Pipe(append([]OperatorFunc{o}, fns...)...)
}

func Pipe(fns ...OperatorFunc) OperatorFunc {
	return func(observable Observable) Observable {
		for _, fn := range fns {
			observable = fn(observable)
		}
		return observable
	}
}

func Range(start, end int) Observable {
	return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
		go func() {
			for i := start; i < end; i++ {
				v <- i
			}
			c <- true
		}()
		return func() {
			c.Complete()
		}
	})
}

func Of(values ...Value) Observable {
	return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
		go func() {
			for _, value := range values {
				v.Next(value)
			}
			c.Complete()
		}()
		return func() {}
	})
}
