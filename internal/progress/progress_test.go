package progress

import (
	"errors"
	"fmt"
	"reflect"
	"sync/atomic"
	"testing"
)

/* -------------------------------- Test: New ------------------------------- */

func TestNew(t *testing.T) {
	tests := []struct {
		size uint64
		want Progress
		err  error
	}{
		// Invalid inputs
		{size: 0, want: Progress{}, err: ErrInvalidTotal},

		// Valid inputs
		{size: 10, want: Progress{total: float64(10), current: &atomic.Uint64{}}},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			// When: A new 'Progress' struct is created with the specified size.
			got, err := New(tc.size)

			// Then: It matches the expected value.
			if !errors.Is(err, tc.err) {
				t.Fatalf("err: got %#v, want %#v", err, tc.err)

			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}

/* ------------------------ Test: Progress.Percentage ----------------------- */

func TestProgressPercentage(t *testing.T) {
	tests := []struct {
		current, size uint64
		want          float64
		err           error
	}{
		// Invalid inputs
		{size: 0, want: 0, err: ErrInvalidTotal},

		// Valid inputs
		{current: 0, size: 10, want: 0},
		{current: 5, size: 10, want: 0.5},
		{current: 10, size: 10, want: 1},
		{current: 15, size: 10, want: 1.5},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			// Given: A 'Progress' struct with the specified size.
			p := Progress{total: float64(tc.size), current: &atomic.Uint64{}}

			// Given: The specified progress is already made.
			p.add(uint64(tc.current))

			// When: The current progress percentage is collected.
			got := p.Percentage()

			// Then: It matches the expected value of 'current' / 'total'.
			if got != tc.want {
				t.Fatalf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}

/* --------------------------- Test: Progress.add --------------------------- */

func TestProgressAdd(t *testing.T) {
	tests := []struct {
		size, add  uint64
		want       uint64
		percentage float64
	}{
		{size: 10, add: 0, want: 0, percentage: 0.0},
		{size: 10, add: 5, want: 5, percentage: 0.5},
		{size: 10, add: 10, want: 10, percentage: 1.0},
		{size: 10, add: 20, want: 20, percentage: 2.0},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			// Given: A 'Progress' struct with the specified size.
			p, err := New(tc.size)
			if !errors.Is(err, nil) {
				t.Fatalf("err: got %#v, want %#v", err, nil)

			}

			// When: The specified progress amount is added.
			got := p.add(uint64(tc.add))

			// Then: It returns the expected new value.
			if got != tc.want {
				t.Fatalf("output: got %#v, want %#v", got, tc.want)
			}

			// Then: The reported progress reflects the added value.
			if got := p.Percentage(); got != tc.percentage {
				t.Fatalf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}
