package rxgo

import (
	"testing"
	"time"
)

func TestBackoff(t *testing.T) {
	Range(0, 5).
		Pipe(
			Backoff(10*time.Millisecond),
			Print(10*time.Millisecond),
		).
		Subscribe(OnNext(func(v Value) {
			expected := "-0--1----2--------3----------------4c"
			if !comparePrintedOutput(v, expected) {
				t.Errorf("Test failed, expected %v, got %v", expected, v)
			}
			t.Log(v)
		})).
		Wait()
}

func TestBackoffWithMax(t *testing.T) {
	Range(0, 5).
		Pipe(
			BackoffWithMax(10*time.Millisecond, 50*time.Millisecond),
			Print(10*time.Millisecond),
		).
		Subscribe(OnNext(func(v Value) {
			expected := "-0--1----2-----3-----4c"
			if !comparePrintedOutput(v, expected) {
				t.Errorf("Test failed, expected %v, got %v", expected, v)
			}
			t.Log(v)
		})).
		Wait()
}

func comparePrintedOutput(v Value, expected string) bool {
	return v.(string) == expected
}
