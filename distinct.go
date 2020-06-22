package rxgo

import (
	"fmt"
	"reflect"
)

func DistinctUntilChanged(compare func(a, b Value) bool) OperatorFunc {
	var prevVal Value
	return func(o Observable) Observable {
		return o.Pipe(
			Filter(func(val Value) bool {
				if compare(prevVal, val) {
					return false
				}
				prevVal = val
				return true
			}),
		)
	}
}

func Distinct(keyFn func(v Value) string) OperatorFunc {
	distinctMap := make(map[string]struct{})
	return Filter(func(val Value) bool {
		key := keyFn(val)
		if _, ok := distinctMap[key]; !ok {
			distinctMap[key] = struct{}{}
			return true
		}
		return false
	})
}

func DeepEqual(a, b Value) bool {
	return reflect.DeepEqual(a, b)
}

func Sprint(a Value) string {
	return fmt.Sprint(a)
}
