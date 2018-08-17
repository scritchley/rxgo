package rxgo

import (
	"time"
)

// Interval returns IntervalWithScheduler using an AsyncScheduler.
func Interval(d time.Duration) Observable {
	return IntervalWithScheduler(d, AsyncScheduler{})
}

func IntervalWithScheduler(d time.Duration, scheduler Scheduler) Observable {
	return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
		go func() {
			count := 0
			for {
				scheduler.Schedule(func() {
					v.Next(count)
					count++
				}, d)
			}
		}()
		return c.Complete
	})
}
