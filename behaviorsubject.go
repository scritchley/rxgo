package rxgo

import (
	"sync"
)

type BehaviorSubject struct {
	mtx sync.Mutex
	Value
	Subject
}

func NewBehaviorSubject(initialValue Value) Subject {
	return &BehaviorSubject{Value: initialValue, Subject: NewSubject()}
}

func (b *BehaviorSubject) Pipe(fns ...OperatorFunc) Observable {
	return Pipe(fns...)(b)
}

func (b *BehaviorSubject) Next(v Value) {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	b.Value = v
	b.Subject.Next(v)
}

func (b *BehaviorSubject) Subscribe(obs Observer) Subscription {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	go obs.Next(b.Value)
	return b.Subject.Subscribe(obs)
}
