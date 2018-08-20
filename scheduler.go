package rxgo

import (
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
	go work()
}
