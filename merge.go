package rxgo

import (
	"sync"
)

func Merge(observables ...Observable) Observable {
	return Create(func(v ValueChan, e ErrChan, c CompleteChan) TeardownFunc {
		psub := NewSubscription(nil)
		var wg sync.WaitGroup
		wg.Add(len(observables))
		for _, observable := range observables {
			psub.Add(observable.Subscribe(
				OnNext(v.Next).
					OnErr(e.Error).
					OnComplete(func() {
						wg.Done()
					}),
			).Unsubscribe)
		}
		go func() {
			wg.Wait()
			c <- true
		}()
		return psub.Unsubscribe
	})
}
