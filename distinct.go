package rxgo

import "reflect"

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

func DeepEqual(a, b Value) bool {
	return reflect.DeepEqual(a, b)
}
