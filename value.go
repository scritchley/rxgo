package rxgo

type Value interface{}

func IsValueSlice(v Value) bool {
	_, ok := v.([]Value)
	return ok
}

func ToValueSlice(v Value) []Value {
	if vals, ok := v.([]Value); ok {
		return vals
	}
	return nil
}
