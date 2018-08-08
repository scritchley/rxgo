package rxgo

import (
	"errors"
	"testing"
)

func TestRetry(t *testing.T) {

	Range(0, 10).
		Pipe(
			Delay(tenMilliseconds),
			ConcatMap(func(val Value) Observable {
				if val.(int) > 5 {
					return Throw(errors.New("error"))
				}
				return Of(val)
			}),
			Retry(1),
			Print(tenMilliseconds),
		).
		Subscribe(OnNext(func(v Value) {
			expected := "-0-1-2-3-4-5--0-1-2-3-4-5-e"
			if !comparePrintedOutput(v, expected) {
				t.Errorf("Test failed, expected %v, got %v", expected, v)
			}
			t.Log(v)
		})).
		Wait()

}
