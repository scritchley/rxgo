package rxgo

import (
	"testing"
	"time"
)

func TestTake(t *testing.T) {

	printer := NewPrinter(time.Millisecond)

	Interval(time.Millisecond).
		Pipe(
			Take(10),
		).
		Subscribe(printer).Wait()

	t.Log(printer.String())

}

func TestTakeUntil(t *testing.T) {

	printer := NewPrinter(time.Millisecond)

	Interval(5 * time.Millisecond).
		Pipe(
			TakeUntil(
				Interval(7 * time.Millisecond),
			),
		).
		Subscribe(printer).Wait()

	t.Log(printer.String())

}
