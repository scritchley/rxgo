package rxgo

// func Race(observables ...Observable) Observable {
// 	Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
// 		firstEmit := make(ValueChan, 1)
// 		for _, observable := range observables {
// 			observable.Subscribe(OnNext(func(v Value) {
// 				v <- firstEmit
// 			}))
// 		}

// 		return ob

// 	})
// }
