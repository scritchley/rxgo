package rxgo

type observable struct {
	Subscriber
}

func Create(subscriber Subscriber) Observable {
	return observable{
		subscriber,
	}
}

func (c observable) Subscribe(obs Observer) Subscription {
	valueCh := make(ValueChan)
	errCh := make(ErrChan, 1)
	completeCh := make(CompleteChan, 1)
	td := c.Subscriber(valueCh, errCh, completeCh)
	sub := NewSubscription(td)
	go func() {
		defer sub.Unsubscribe()
	LOOP:
		for {
			select {
			case v := <-valueCh:
				if obs != nil {
					obs.Next(v)
				}
			case e := <-errCh:
				if obs != nil {
					obs.Err(e)
				}
				break LOOP
			case <-completeCh:
				if obs != nil {
					obs.Complete()
				}
				break LOOP
			}
		}
	}()
	return sub
}

func (c observable) Pipe(ops ...OperatorFunc) Observable {
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

type Observables []Observable

func (o Observables) Map(fn func(Observable) Observable) Observables {
	mapped := make([]Observable, len(o))
	for i, observable := range o {
		mapped[i] = fn(observable)
	}
	return mapped
}

type Subscribable interface {
	Subscribe(Observer) Subscription
}

type Observer interface {
	NextObserver
	ErrObserver
	CompletionObserver
}

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
	select {
	case c <- true:
	default:
	}
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
			for i := start; i <= end; i++ {
				v <- i
			}
			c <- true
		}()
		return c.Complete
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
		return nil
	})
}

type noOpObserver struct{}

func (noOpObserver) Next(Value) {}

func (noOpObserver) Err(error) {}

func (noOpObserver) Complete() {}
