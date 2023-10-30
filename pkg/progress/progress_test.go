package progress

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

/* ------------------------- Test: TestNewWithTotal ------------------------- */

func TestNewWithTotal(t *testing.T) {
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
			got, err := NewWithTotal(tc.size)

			// Then: It matches the expected value.
			if !errors.Is(err, tc.err) {
				t.Errorf("err: got %#v, want %#v", err, tc.err)

			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("output: got %#v, want %#v", got, want)
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
			p, err := NewWithTotal(tc.size)
			if !errors.Is(err, nil) {
				t.Errorf("err: got %#v, want %#v", err, nil)

			}

			// Given: The specified progress is already made.
			p.Add(uint64(tc.current))

			// When: The current progress percentage is collected.
			// Then: It matches the expected value of 'current' / 'total'.
			if got := p.Percentage(); got != tc.want {
				t.Errorf("output: got %#v, want %#v", got, tc.want)
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
			p.total.Store(tc.total)

			// Given: The specified progress is already made.
			p.Add(tc.current)

			// When: The current progress is reset.
			p.Reset()

			// Then: The 'current' progress is reset to '0'.
			if got := p.Current(); got != 0 {
				t.Errorf("output: got %#v, want %#v", got, 0)
			}

			// Then: The 'total' value is unchanged.
			if got := p.total.Load(); got != tc.total {
				t.Errorf("output: got %#v, want %#v", got, tc.total)
			}
		})
	}
}

/* ------------------------- Test: Progress.SetTotal ------------------------ */

func TestProgressSetTotal(t *testing.T) {
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
			p.Add(tc.current)

			// When: The 'total' is set to the new value.
			err := p.SetTotal(tc.next)

			// Then: An error is returned if the new value is invalid.
			if !errors.Is(err, tc.err) {
				t.Errorf("err: got %#v, want %#v", err, tc.err)
			}

			// Then: The 'current' progress is unchanged.
			if got := p.current.Load(); got != tc.current {
				t.Errorf("output: got %#v, want %#v", got, tc.current)
			}

			// Then: The 'total' field is updated to the desired value.
			if got := p.Total(); got != tc.want {
				t.Errorf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}
