package rxgo

import (
	"reflect"
	"sync"
)

type Subscription interface {
	Unsubscribe()
	Add(TeardownFunc) Subscription
	Remove(Subscription)
	Wait()
}

type subscription struct {
	mtx sync.Mutex
	TeardownFunc
	done          chan bool
	subscriptions []Subscription
}

func NewSubscription(teardown TeardownFunc) *subscription {
	return &subscription{TeardownFunc: teardown, done: make(chan bool, 1)}
}

func (s *subscription) Unsubscribe() {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	select {
	case s.done <- true:
		for _, sub := range s.subscriptions {
			sub.Unsubscribe()
		}
		if s.TeardownFunc != nil {
			s.TeardownFunc()
		}
	default:
	}
}

func (s *subscription) Add(teardown TeardownFunc) Subscription {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	sub := NewSubscription(teardown)
	s.subscriptions = append(s.subscriptions, sub)
	return sub
}

func (s *subscription) Remove(subscription Subscription) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	for i := range s.subscriptions {
		if reflect.DeepEqual(s.subscriptions[i], subscription) {
			s.subscriptions = append(s.subscriptions[:i], s.subscriptions[i+1:]...)
		}
	}
}

func (s *subscription) Wait() {
	<-s.done
}
