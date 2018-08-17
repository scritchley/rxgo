package rxgo

import (
	"fmt"
	"reflect"
)

func Assert(vals ...Value) OperatorFunc {
	return Pipe(
		ToArray(),
		ConcatMap(func(val Value) Observable {
			if !reflect.DeepEqual(vals, val) {
				return Throw(fmt.Errorf("%v != %v", vals, val))
			}
			return Of(val)
		}),
	)
}
