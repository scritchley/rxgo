package rxgo

import (
	"sync"
)

// ForkJoin waits for all observables to complete and will then emit the last emitted value from each as an array.
func ForkJoin(observables ...Observable) Observable {
	return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
		finalValues := make([]Value, len(observables))
		innerSub := NewSubscription(nil)
		var wg sync.WaitGroup
		for i, observable := range observables {
			wg.Add(1)
			index := i
			innerSub.Add(observable.Pipe(Last()).Subscribe(OnNext(func(v Value) {
				finalValues[index] = v
			}).OnErr(e.Error).OnComplete(wg.Done)).Unsubscribe)
		}
		go func() {
			wg.Wait()
			v.Next(finalValues)
			c.Complete()
		}()
		return call(
			innerSub.Unsubscribe,
		)
	})
}
