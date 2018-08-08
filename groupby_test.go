package rxgo

import (
	"testing"
)

type m map[string]string

func TestGroupBy(t *testing.T) {

	Merge(
		Range(0, 10),
		Range(0, 10),
		Range(0, 10),
	).
		Pipe(
			Delay(tenMilliseconds),
			GroupBy(
				func(v Value) Value {
					return v
				},
				ToArray(),
			),
			Print(tenMilliseconds),
		).
		Subscribe(OnNext(func(v Value) {
			t.Log(v)
		})).
		Wait()

}
