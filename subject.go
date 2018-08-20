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
	index     int
	observers map[int]Observer
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
	return &subject{observers: make(map[int]Observer)}
}

func (s *subject) Subscribe(obs Observer) Subscription {
	s.mtx.Lock()
	index := s.index
	s.index++
	s.observers[index] = obs
	s.mtx.Unlock()
	return NewSubscription(func() {
		s.mtx.Lock()
		defer s.mtx.Unlock()
		delete(s.observers, index)
	})
}

type replaySubject struct {
	mtx    sync.Mutex
	buffer []Value
	Subject
}

func NewReplaySubject(buffer int) Subject {
	return &replaySubject{buffer: make([]Value, buffer), Subject: NewSubject()}
}

func (r *replaySubject) Next(v Value) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.buffer = append(r.buffer[1:], v)
	r.Subject.Next(v)
}

func (r *replaySubject) Subscribe(obs Observer) Subscription {
	return Create(
		func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
			r.mtx.Lock()
			for _, val := range r.buffer {
				v.Next(val)
			}
			r.mtx.Unlock()
			println("sub")
			return r.Subject.Subscribe(
				OnNext(v.Next).OnErr(e.Error).OnComplete(c.Complete),
			).Unsubscribe
		},
	).Subscribe(obs)
}
