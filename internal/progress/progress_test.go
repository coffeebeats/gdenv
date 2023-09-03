package progress

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

/* -------------------------------- Test: New ------------------------------- */

func TestNew(t *testing.T) {
	tests := []struct {
		size uint64
		want *Progress
		err  error
	}{
		// Invalid inputs
		{size: 0, want: &Progress{}, err: ErrInvalidTotal},

		// Valid inputs
		{size: 10},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			var want *Progress

			// Given: The correctly instantiated expected struct.
			if tc.size > 0 {
				want = &Progress{}
				want.total.Store(tc.size)
			}

			// When: A new 'Progress' struct is created with the specified size.
			got, err := New(tc.size)

			// Then: It matches the expected value.
			if !errors.Is(err, tc.err) {
				t.Fatalf("err: got %#v, want %#v", err, tc.err)

			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("output: got %#v, want %#v", got, want)
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
		{current: 0, size: 10, want: 0},
		{current: 5, size: 10, want: 0.5},
		{current: 10, size: 10, want: 1},
		{current: 15, size: 10, want: 1.5},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			// Given: A 'Progress' struct with the specified size.
			p, err := New(tc.size)
			if !errors.Is(err, nil) {
				t.Fatalf("err: got %#v, want %#v", err, nil)

			}

			// Given: The specified progress is already made.
			p.add(uint64(tc.current))

			// When: The current progress percentage is collected.
			// Then: It matches the expected value of 'current' / 'total'.
			if got := p.Percentage(); got != tc.want {
				t.Fatalf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}

/* -------------------------- Test: Progress.Reset -------------------------- */

func TestProgressReset(t *testing.T) {
	tests := []struct {
		current, total uint64
	}{
		{current: 0, total: 1},
		{current: 1, total: 0},
		{current: 1, total: 1},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%d-%d", tc.current, tc.total), func(t *testing.T) {
			// Given: A 'Progress' struct with the specified total.
			p := Progress{}
			p.current.Store(tc.current)
			p.total.Store(tc.total)

			// Given: The specified progress is already made.
			p.add(tc.current)

			// When: The current progress is reset.
			p.Reset()

			// Then: The 'current' progress is reset to '0'.
			if got := p.current.Load(); got != 0 {
				t.Fatalf("output: got %#v, want %#v", got, 0)
			}

			// Then: The 'total' value is unchanged.
			if got := p.total.Load(); got != tc.total {
				t.Fatalf("output: got %#v, want %#v", got, tc.total)
			}
		})
	}
}

/* -------------------------- Test: Progress.Total -------------------------- */

func TestProgressTotal(t *testing.T) {
	tests := []struct {
		current, total, next, want uint64
		err                        error
	}{
		// Invalid inputs
		{current: 0, total: 1, next: 0, want: 1, err: ErrInvalidTotal},
		{current: 1, total: 0, next: 0, want: 0, err: ErrInvalidTotal},

		// Valid inputs
		{current: 0, total: 1, next: 10, want: 10},
		{current: 1, total: 0, next: 10, want: 10},
		{current: 1, total: 1, next: 10, want: 10},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%d-%d", tc.current, tc.total), func(t *testing.T) {
			// Given: A 'Progress' struct with the specified total.
			p := Progress{}
			p.total.Store(tc.total)

			// Given: The specified progress is already made.
			p.add(tc.current)

			// When: The 'total' is set to the new value.
			err := p.Total(tc.next)

			// Then: An error is returned if the new value is invalid.
			if !errors.Is(err, tc.err) {
				t.Fatalf("err: got %#v, want %#v", err, tc.err)
			}

			// Then: The 'current' progress is unchanged.
			if got := p.current.Load(); got != tc.current {
				t.Fatalf("output: got %#v, want %#v", got, tc.current)
			}

			// Then: The 'total' field is updated to the desired value.
			if got := p.total.Load(); got != tc.want {
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
			// Then: It returns the expected new value.
			if got := p.add(uint64(tc.add)); got != tc.want {
				t.Fatalf("output: got %#v, want %#v", got, tc.want)
			}

			// Then: The reported progress reflects the added value.
			if got := p.Percentage(); got != tc.percentage {
				t.Fatalf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}
