package rxgo

import (
	"time"
)

func Interval(d time.Duration) Observable {
	return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
		go func() {
			count := 0
			for {
				time.Sleep(d)
				v.Next(count)
				count++
			}
		}()
		return c.Complete
	})
}
