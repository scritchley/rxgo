package rxgo

import (
	"sync"
)

type behaviorSubject struct {
	mtx sync.Mutex
	Value
	Subject
}

func NewBehaviorSubject(initialValue Value) Subject {
	return &behaviorSubject{Value: initialValue, Subject: NewSubject()}
}

func (b *behaviorSubject) Pipe(fns ...OperatorFunc) Observable {
	return Pipe(fns...)(b)
}

func (b *behaviorSubject) Next(v Value) {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	b.Value = v
	b.Subject.Next(v)
}

func (b *behaviorSubject) Subscribe(obs Observer) Subscription {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	go obs.Next(b.Value)
	return b.Subject.Subscribe(obs)
}
