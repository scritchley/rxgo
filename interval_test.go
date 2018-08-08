package rxgo

import (
	"testing"
	"time"
)

func TestInterval(t *testing.T) {

	printer := NewPrinter(time.Millisecond)

	Interval(time.Millisecond).
		Pipe(
			Take(10),
		).
		Subscribe(printer).Wait()

	t.Log(printer.String())

}
