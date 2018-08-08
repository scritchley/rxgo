package rxgo

import (
	"errors"
	"testing"
	"time"
)

func TestThrow(t *testing.T) {

	printer := NewPrinter(time.Millisecond)

	Throw(errors.New("test error")).
		Subscribe(printer).Wait()

	t.Log(printer.String())

}
