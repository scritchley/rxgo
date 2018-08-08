package rxgo

import "sync/atomic"

type Subject struct {
	Observable
	Observer
}

func NewSubject() Subject {
	obs := NewChanObserver()
	return Subject{
		Observer:   obs,
		Observable: obs.From(),
	}
}

type BehaviorSubject struct {
	Subject
	atomic.Value
}

func (b *BehaviorSubject) GetValue() Value {
	return b.Value.Load()
}

func (b *BehaviorSubject) Next(v Value) {
	b.Store(v)
	b.Subject.Next(v)
}
