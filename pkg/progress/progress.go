package progress

import (
	"errors"
	"fmt"
	"sync/atomic"
)

var ErrInvalidTotal = errors.New("invalid total")

/* -------------------------------------------------------------------------- */
/*                              Struct: Progress                              */
/* -------------------------------------------------------------------------- */

// A thread-safe progress tracker; reports percentages exactly (i.e. greater
// than '1.0' is possible), but will panic on denominators <= '0'.
//
// NOTE: A 'Progress' struct must not be copied after first use.
type Progress struct {
	current, total atomic.Uint64
}

/* ------------------------- Function: NewWithTotal ------------------------- */

// Creates a new 'Progress' struct with the specified 'total' size.
func NewWithTotal(total uint64) (*Progress, error) {
	var progress Progress

	if err := progress.SetTotal(total); err != nil {
		return nil, err
	}

	return &progress, nil
}

/* ------------------------------- Method: Add ------------------------------ */

// Adds the specified amount to the current progress and returns the new
// 'current' value. This method is thread-safe.
func (p *Progress) Add(n uint64) uint64 {
	return p.current.Add(n)
}

/* ----------------------------- Method: Current ---------------------------- */

// Return the 'current' value for the 'Progress' struct. This method is
// thread-safe.
func (p *Progress) Current() uint64 {
	return p.current.Load()
}

/* --------------------------- Method: Percentage --------------------------- */

// Retrieves the current progress as a decimal fraction. This method is
// thread-safe.
//
// NOTE: This method returns '0.0' if the underlying total is unset.
func (p *Progress) Percentage() float64 {
	total := p.total.Load()
	if total == 0 {
		return 0
	}

	return float64(p.current.Load()) / float64(total)
}

/* --------------------------- Method: Reset --------------------------- */

// Resets the current progress to '0'. This method is thread-safe.
//
// NOTE: This method does *NOT* modify the 'total' value.
func (p *Progress) Reset() {
	p.current.Store(0)
}

/* ---------------------------- Method: SetTotal ---------------------------- */

// Modifies the 'total' size value for the 'Progress' struct. This method is
// thread-safe.
func (p *Progress) SetTotal(total uint64) error {
	if total == 0 {
		return fmt.Errorf("%w: %d", ErrInvalidTotal, total)
	}

	p.total.Store(total)

	return nil
}

/* ------------------------------ Method: Total ----------------------------- */

// Return the 'total' value for the 'Progress' struct. This method is
// thread-safe.
func (p *Progress) Total() uint64 {
	return p.total.Load()
}
