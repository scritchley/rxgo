package rxgo

import "time"

type Scheduler interface {
	Schedule(func())
	Stop()
}

type TickerScheduler struct {
	*time.Ticker
}

func (t TickerScheduler) Schedule(fn func() bool) {
	go func() {
		for range t.C {
			if !fn() {
				t.Stop()
			}
		}
	}()
}

func (t TickerScheduler) Stop() {
	t.Stop()
}

func NewTickerScheduler(d time.Duration) TickerScheduler {
	return TickerScheduler{
		Ticker: time.NewTicker(d),
	}
}
