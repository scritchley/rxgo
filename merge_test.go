package rxgo

import (
	"testing"
)

func TestMerge(t *testing.T) {

	Merge(
		Of(0),
		Of(1).Pipe(Delay(10*tenMilliseconds)),
		Of(2).Pipe(Delay(20*tenMilliseconds)),
	).
		Pipe(
			Print(tenMilliseconds),
		).
		Subscribe(OnNext(func(v Value) {
			expected := "0----------1----------2c"
			if !comparePrintedOutput(v, expected) {
				t.Errorf("Test failed, expected %v, got %v", expected, v)
			}
			t.Log(v)
		})).
		Wait()

}
