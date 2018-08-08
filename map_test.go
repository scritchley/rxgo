package rxgo

import (
	"testing"
	"time"
)

func TestConcatMap(t *testing.T) {

	Interval(tenMilliseconds).Pipe(
		Take(10),
		ConcatMap(func(v Value) Observable {
			return Of(1, 2, 3).Pipe(
				Delay(tenMilliseconds),
			)
		}),
		Print(tenMilliseconds),
	).Subscribe(OnNext(func(v Value) {
		t.Log(v)
	})).Wait()

}

func TestMergeMap(t *testing.T) {

	printer := NewPrinter(time.Millisecond)

	Interval(1*time.Millisecond).Pipe(
		Take(10),
		MergeMap(func(v Value) Observable {
			return Of(1, 2)
		}),
	).Subscribe(printer).Wait()

	t.Log(printer.String())

}
