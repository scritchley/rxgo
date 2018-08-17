package rxgo

import (
	"sync"
)

type Subject interface {
	Observable
	Observer
}

type subject struct {
	mtx       sync.Mutex
	observers []Observer
}

func (s *subject) Next(v Value) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	for _, observer := range s.observers {
		observer.Next(v)
	}
}

func (s *subject) Err(e error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	for _, observer := range s.observers {
		observer.Err(e)
	}
}

func (s *subject) Complete() {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	for _, observer := range s.observers {
		observer.Complete()
	}
}

func (s *subject) Pipe(fns ...OperatorFunc) Observable {
	return Pipe(fns...)(s)
}

func NewSubject() Subject {
	return &subject{}
}

func (s *subject) Subscribe(obs Observer) Subscription {
	s.mtx.Lock()
	s.observers = append(s.observers, obs)
	l := len(s.observers)
	s.mtx.Unlock()
	return NewSubscription(func() {
		s.mtx.Lock()
		defer s.mtx.Unlock()
		s.observers = append(s.observers[:l], s.observers[l+1:]...)
	})
}
