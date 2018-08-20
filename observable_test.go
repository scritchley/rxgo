package rxgo

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

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
		"Merge should merge the output of the provided observables",
		Merge(
			Of(1),
			Of(1),
			Of(1),
		).
			Pipe(
				Assert(1, 1, 1),
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
		"ToArray should append all emitted values to slice and emit once",
		Of(1, 2, 3).
			Pipe(
				ToArray(),
				Assert([]Value{1, 2, 3}),
			),
	},
	{
		"Last should emit the last value",
		Of(1, 2, 3).
			Pipe(
				Last(),
				Assert(3),
			),
	},
	{
		"Reduce should apply reduce func to values",
		Range(0, 5).
			Pipe(
				Reduce(func(acc, value Value) Value {
					return acc.(int) + value.(int)
				}, nil),
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
				ConcatMapTo(
					Of(1, 2, 3),
				),
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
					return Of(value).Pipe(
						Delay(time.Duration(value.(int)) * time.Millisecond),
					)
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
		Interval(40*time.Millisecond).
			Pipe(
				TakeUntil(
					Interval(100*time.Millisecond),
				),
				Assert(0, 1),
			),
	},
	{
		"ForkJoin should return an array containing the results in order",
		ForkJoin(
			Of(0, 1),
			Of(0, 2),
			Of(0, 3),
			Of(0, 4),
		).Pipe(
			Assert([]Value{1, 2, 3, 4}),
		),
	},
	{
		"GroupBy should return an array containing",
		Of(1, 1, 1).
			Pipe(
				GroupBy(func(v Value) Value {
					return v
				},
					ToArray(),
				),
			).
			Pipe(
				Assert([]Value{1, 1, 1}),
			),
	},
	{
		"Share should share source between multiple subscribers",
		Of(1, 2, 3).
			Pipe(
				Share(),
				Assert(1, 2, 3),
			),
	},
	{
		"IgnoreElements should ignore all elements",
		Of(1, 2, 3).
			Pipe(
				IgnoreElements(),
				Assert(),
			),
	},
	{
		"DistinctUntilChanged should only emit when the current value is different than the last.",
		Of(1, 1, 2, 2, 3, 4, 5, 5).
			Pipe(
				DistinctUntilChanged(DeepEqual),
				Assert(1, 2, 3, 4, 5),
			),
	},
	{
		"Catch should catch an error",
		Of(1).
			Pipe(
				ConcatMap(func(val Value) Observable {
					return Throw(errors.New("err"))
				}),
				Catch(func(e error) (Value, error) {
					if reflect.DeepEqual(e, errors.New("err")) {
						return e, nil
					}
					return nil, e
				}),
				Assert(errors.New("err")),
			),
	},
	{
		"Retry should retry the observable on receiving an error",
		Of(1, 2, 3).
			Pipe(
				ConcatMap(func(val Value) Observable {
					if val.(int) > 2 {
						return Throw(errors.New("err"))
					}
					return Of(val)
				}),
				Retry(2),
				Catch(func(e error) (Value, error) {
					if reflect.DeepEqual(e, errors.New("err")) {
						return e, nil
					}
					return nil, e
				}),
				Assert(1, 2, 1, 2, 1, 2, errors.New("err")),
			),
	},
	{
		"ExhaustMap should map to inner observable, ignore other values until that observable completes.",
		Interval(time.Millisecond).
			Pipe(
				Take(2),
				ExhaustMap(func(val Value) Observable {
					return Of(1, 2, 3).Pipe(
						Delay(10 * time.Millisecond),
					)
				}),
				Assert(1, 2, 3),
			),
	},
	{
		"Timeout should timeout if the observable does not emit within the provided timeout duration.",
		Interval(time.Second).
			Pipe(
				Take(3),
				Timeout(time.Millisecond),
				Assert(),
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
