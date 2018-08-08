package rxgo

import (
	"testing"
	"time"
)

func TestWindow(t *testing.T) {

	printer := NewPrinter(time.Millisecond)

	Interval(time.Millisecond).Pipe(
		Take(10),
		Window(Interval(20*time.Millisecond)),
	).Subscribe(printer).Wait()

	t.Log(printer.String())

}
