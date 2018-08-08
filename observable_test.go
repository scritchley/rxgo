package rxgo

import "time"

var tenMilliseconds = 10 * time.Millisecond

// func TestObservable(t *testing.T) {

// 	// Merge(
// 	// 	Range(0, 10),
// 	// 	Range(10, 20),
// 	// 	Range(20, 30),
// 	// ).Subscribe(
// 	// 	OnNext(func(v Value) {
// 	// 		fmt.Println(v)
// 	// 	}).OnErr(func(err error) {
// 	// 		fmt.Println("error", err)
// 	// 	}).OnComplete(func() {
// 	// 		fmt.Println("completed")
// 	// 	}),
// 	// ).Wait()

// 	// Range(0, 10).Pipe(
// 	// 	PairWise(),
// 	// 	Delay(time.Second),
// 	// 	Do(func() {
// 	// 		println("hi")
// 	// 	}),
// 	// ThrottleTime(200*time.Millisecond),
// 	// MergeMap(func(v Value) Observable {
// 	// 	return Range(0, 10)
// 	// }),
// 	// Buffer(10),
// 	// Map(func(v Value) Value {
// 	// 	return true
// 	// }),

// 	// Create(func(obs Observer) TeardownFunc {
// 	// 	for i := 0; i < 10; i++ {
// 	// 		obs.Next(i)
// 	// 		time.Sleep(time.Second)
// 	// 	}
// 	// 	return func() {
// 	// 		t.Log("done")
// 	// 	}
// 	// }).Pipe(
// 	// 	Take(3),
// 	// ).Subscribe(
// 	// 	NewChanObserver(
// 	// 		OnNext(func(v Value) {
// 	// 			fmt.Println(v)
// 	// 		}).OnComplete(func() {
// 	// 			fmt.Println("completed")
// 	// 		}).OnErr(func(err error) {
// 	// 			fmt.Println("error", err)
// 	// 		}),
// 	// 	),
// 	// )

// }
