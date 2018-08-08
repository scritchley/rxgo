package rxgo

import (
	"testing"
)

func TestReduce(t *testing.T) {

	Range(0, 5).Pipe(
		Reduce(func(acc, value Value) Value {
			return acc.(int) + value.(int)
		}),
		Print(tenMilliseconds),
	).Subscribe(OnNext(func(v Value) {
		expected := "10c"
		if !comparePrintedOutput(v, expected) {
			t.Errorf("Test failed, expected %v, got %v", expected, v)
		}
		t.Log(v)
	})).Wait()

}
