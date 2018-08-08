package rxgo

import (
	"testing"
)

func TestZip(t *testing.T) {

	Zip(
		Interval(tenMilliseconds),
		Interval(tenMilliseconds),
		// Range(10, 20),
		// Range(10, 20),
	).Pipe(
		Take(10),
		Print(tenMilliseconds),
	).Subscribe(OnNext(func(v Value) {
		t.Log(v)
	})).Wait()

}
