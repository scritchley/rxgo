package rxgo

import (
	"testing"
	"time"
)

var tenMilliseconds = 10 * time.Millisecond

var testCases = []struct {
	name string
	Observable
}{
	{
		"Of should emit the provided values",
		Of(1, 2, 3).
			Pipe(
				Assert(1, 2, 3),
			),
	},
	{
		"Range should emit the provided range",
		Range(0, 5).
			Pipe(
				Assert(0, 1, 2, 3, 4, 5),
			),
	},
	{
		"Interval should emit sequential numbers",
		Interval(time.Millisecond).
			Pipe(
				Take(5),
				Assert(0, 1, 2, 3, 4),
			),
	},
	{
		"Pipe should pipe via the provided operators",
		Of(1, 2, 3).
			Pipe(
				Pipe(
					Assert(1, 2, 3),
				),
			),
	},
	{
		"Reduce should apply reduce func to values",
		Range(0, 5).
			Pipe(
				Reduce(func(acc, value Value) Value {
					return acc.(int) + value.(int)
				}),
				Assert(15),
			),
	},
	{
		"Map should apply map func to values",
		Range(0, 5).
			Pipe(
				Map(func(value Value) Value {
					return value.(int) * 2
				}),
				Assert(0, 2, 4, 6, 8, 10),
			),
	},
	{
		"MapTo should map all values to the provided value",
		Range(0, 5).
			Pipe(
				MapTo(1),
				Assert(1, 1, 1, 1, 1, 1),
			),
	},
	{
		"ConcatMap should subscribe to inner observables sequentially",
		Of(1, 2, 3).
			Pipe(
				ConcatMap(func(value Value) Observable {
					return Of(value, value, value)
				}),
				Assert(1, 1, 1, 2, 2, 2, 3, 3, 3),
			),
	},
	{
		"ConcatMapTo should subscribe to the provided observable",
		Of(1, 2, 3).
			Pipe(
				ConcatMapTo(Of(1, 2, 3)),
				Assert(1, 2, 3, 1, 2, 3, 1, 2, 3),
			),
	},
	{
		// The effect of this test should be to reverse the values provided
		// to Of, as each is delayed according to it's value in milliseconds.
		"MergeMap should subscribe to inner observables in parallel",
		Of(300, 200, 100).
			Pipe(
				MergeMap(func(value Value) Observable {
					return Of(value).Pipe(Delay(time.Duration(value.(int)) * time.Millisecond))
				}),
				Assert(100, 200, 300),
			),
	},
	{
		"SwitchMap should subscribe to inner observable and cancel if pending",
		Range(0, 10).
			Pipe(
				SwitchMap(func(value Value) Observable {
					return Of(value).
						Pipe(
							Delay(300 * time.Millisecond),
						)
				}),
				Assert(10),
			),
	},
	{
		"BehaviorSubject should start with the provided value",
		NewBehaviorSubject("start").
			Pipe(
				Take(1),
				Assert("start"),
			),
	},
	{
		"PairWise should emit the last and current value",
		Range(0, 4).Pipe(
			PairWise(),
			Assert([]Value{0, 1}, []Value{1, 2}, []Value{2, 3}, []Value{3, 4}),
		),
	},
	{
		"Buffer should emit the a slice of values",
		Range(0, 6).Pipe(
			Buffer(3),
			Assert([]Value{0, 1, 2}, []Value{3, 4, 5}),
		),
	},
	{
		"TakeUntil should take values until the provided Observable emits",
		Interval(40*time.Millisecond).Pipe(
			TakeUntil(Interval(100*time.Millisecond)),
			Assert(0, 1),
		),
	},
}

func TestObservables(t *testing.T) {
	for _, tc := range testCases {
		tc.Subscribe(
			OnNext(func(v Value) {
				t.Logf("SUCCESS `%s` -> %v", tc.name, v)
			}).OnErr(func(e error) {
				t.Errorf("FAILED  `%s` -> %v", tc.name, e)
			}),
		).Wait()
	}
}
