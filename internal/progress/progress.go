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

/* ------------------------------ Function: New ----------------------------- */

// Creates a new 'Progress' struct with the specified 'total' size.
func New(total uint64) (*Progress, error) {
	var progress Progress

	if err := progress.Total(total); err != nil {
		return nil, err
	}

	return &progress, nil
}

/* --------------------------- Method: Percentage --------------------------- */

// Retrieves the current progress as a decimal fraction. This method is
// thread-safe.
func (p *Progress) Percentage() float64 {
	return float64(p.current.Load()) / float64(p.total.Load())
}

/* --------------------------- Method: Reset --------------------------- */

// Resets the current progress to '0'. This method is thread-safe.
//
// NOTE: This method does *NOT* modify the 'total' value.
func (p *Progress) Reset() {
	p.current.Store(0)
}

/* ------------------------------ Method: Total ----------------------------- */

// Modifies the 'total' size value for the 'Progress' struct. This method is
// thread-safe.
func (p *Progress) Total(total uint64) error {
	if total == 0 {
		return fmt.Errorf("%w: %d", ErrInvalidTotal, total)
	}

	p.total.Store(total)

	return nil
}

/* ------------------------------- Method: add ------------------------------ */

// Adds the specified amount to the current progress and returns the new
// 'current' value.
func (p *Progress) add(n uint64) uint64 {
	return p.current.Add(n)
}
