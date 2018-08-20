package rxgo

func Merge(observables ...Observable) Observable {
	return Range(0, len(observables)-1).Pipe(
		MergeMap(func(v Value) Observable {
			return observables[v.(int)]
		}),
	)
}
