package rxgo

import (
	"bytes"
	"fmt"
	"math"
	"sync"
	"time"
)

func Print(scheduler Scheduler, resolution time.Duration) OperatorFunc {
	return func(obs Observable) Observable {
		return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
			printer := NewPrinter(scheduler, resolution)
			return obs.Subscribe(OnNext(printer.Next).
				OnErr(func(err error) {
					printer.Err(err)
					v.Next(printer.String())
					c.Complete()
				}).
				OnComplete(func() {
					printer.Complete()
					v.Next(printer.String())
					c.Complete()
				})).Unsubscribe
		})
	}
}

type Printer struct {
	scheduler      Scheduler
	lastUpdateTime time.Time
	mtx            sync.Mutex
	bytes.Buffer
	resolution time.Duration
}

func NewPrinter(scheduler Scheduler, resolution time.Duration) *Printer {
	return &Printer{scheduler: scheduler, resolution: resolution}
}

func (p *Printer) Next(v Value) {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	p.addTicks()
	p.Buffer.WriteString(fmt.Sprint(v))
}

func (p *Printer) Err(err error) {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	p.addTicks()
	p.Buffer.WriteString("e")
}

func (p *Printer) Complete() {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	p.addTicks()
	p.Buffer.WriteString("c")
}

func (p *Printer) addTicks() {
	if !p.lastUpdateTime.IsZero() {
		since := p.scheduler.Now().Sub(p.lastUpdateTime)
		millis := math.Round(float64(since.Nanoseconds() / p.resolution.Nanoseconds()))
		for i := 0; i < int(millis); i++ {
			p.Buffer.WriteString("-")
		}
	}
	p.lastUpdateTime = time.Now()
}
