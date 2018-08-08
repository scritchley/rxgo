package rxgo

import (
	"testing"
	"time"
)

func TestWithLatestFrom(t *testing.T) {

	printer := NewPrinter(time.Millisecond)

	Interval(5*time.Millisecond).
		Pipe(
			WithLatestFrom(
				Interval(10*time.Millisecond),
			),
			Take(1),
		).
		Subscribe(printer.Start()).Wait()

	t.Log(printer.String())

}
