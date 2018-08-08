package rxgo

import (
	"testing"
)

func TestThrottleTime(t *testing.T) {

	Interval(tenMilliseconds).Pipe(
		ThrottleTime(10*tenMilliseconds),
		Take(5),
		Print(tenMilliseconds),
	).Subscribe(OnNext(func(v Value) {
		expected := "-----------0----------10----------20----------30----------40c"
		if !comparePrintedOutput(v, expected) {
			t.Errorf("Test failed, expected %v, got %v", expected, v)
		}
		t.Log(v)
	})).Wait()

}
