# RxGo

This package implements observables for the Go programming language.

# Installation

    go get github.com/scritchley/rxgo

# Getting Started

    Range(0, 5).
        Subscribe(
            OnNext(func(v Value) {
                fmt.Println(v)
            }).
            OnComplete(func(v Value) {
                fmt.Println("completed")
            }),
        ).
        Wait()
    
    // Outputs:
    // 0
    // 1
    // 2
    // 3
    // 4
    // 5
    // completed

## Creating Observables

There are a number of methods for creating source observables, alternatively, you can use a `Subject` to create an observable from any source.

    // Creates an observable that emits an increments int at a fixed duration interval.
    Interval(time.Second)

    // Creates an observable directly from a variable number of values.
    Of(1, 2, 3)

    // Creates an observable that emits integers between the start and end range.
    Range(0, 10)

## Operators

All observables implement the `Pipeable` interface allowing the use of operators. 

eachAs an example, the ConcatMap operator can be used to map each source value to an observable which is merged into the output observable in serial.

    Of(1, 2, 3).
        Pipe(
            ConcatMap(func(v Value) {
                return Of(1,2,3)
            }),
        ).
        Subscribe(
            OnNext(func(v Value) {
                fmt.Println(v)
            }),
        )

    // Output
    // 1
    // 2
    // 3
    // 1
    // 2
    // 3
    // 1
    // 2
    // 3

There are a number of additional operators available that can be used to compose complex observable streams. You can find out more about them at godoc.org/github.com/scritchley/rxgo

## Generating observables for custom types

