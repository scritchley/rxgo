package rxgo

import (
	"sync"
	"time"
)

var (
	DefaultAsyncScheduler Scheduler = AsyncScheduler{}
)

type Scheduler interface {
	Now() time.Time
	Schedule(work func(), delay time.Duration)
}

type AsyncScheduler struct {
}

func (s AsyncScheduler) Now() time.Time {
	return time.Now()
}

func (s AsyncScheduler) Schedule(work func(), delay time.Duration) {
	<-time.After(delay)
	work()
}

type ASAPScheduler struct{}

func (s ASAPScheduler) Now() time.Time {
	return time.Now()
}

func (s ASAPScheduler) Schedule(work func(), delay time.Duration) {
	work()
}

type MockTimeBasedScheduler struct {
	mtx  sync.RWMutex
	base int64
	time.Duration
}

func NewMockTimeBasedScheduler() *MockTimeBasedScheduler {
	return &MockTimeBasedScheduler{base: time.Now().UnixNano()}
}

func (s *MockTimeBasedScheduler) Now() time.Time {
	// s.mtx.Lock()
	// defer s.mtx.Unlock()
	return time.Unix(0, s.base).Add(s.Duration)
}

func (s *MockTimeBasedScheduler) Schedule(work func(), delay time.Duration) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.Duration += delay
	work()
}
