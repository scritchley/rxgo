package rxgo

import (
	"errors"
	"testing"
	"time"
)

func TestCatch(t *testing.T) {

	Range(0, 10).
		Pipe(
			Delay(10*time.Millisecond),
			ConcatMap(func(val Value) Observable {
				if val.(int) > 5 {
					return Throw(errors.New("error"))
				}
				return Of(val)
			}),
			Catch(func(err error) (Value, error) {
				return 6, nil
			}),
			Print(10*time.Millisecond),
		).
		Subscribe(OnNext(func(v Value) {
			t.Log(v)
		})).
		Wait()

}
