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
// than '1.0' is possible), but will fail on denominators <= '0'.
type Progress struct {
	current *atomic.Uint64
	total   float64
}

/* ------------------------------ Function: New ----------------------------- */

// Creates a new 'Progress' struct with the specified 'total' size.
func New(total uint64) (Progress, error) {
	if total == 0 {
		return Progress{}, fmt.Errorf("%w: %d", ErrInvalidTotal, total)
	}

	return Progress{current: &atomic.Uint64{}, total: float64(total)}, nil
}

/* --------------------------- Method: Percentage --------------------------- */

// Retrieves the current progress as a decimal fraction.
func (p *Progress) Percentage() float64 {
	return float64(p.current.Load()) / p.total
}

/* ------------------------------- Method: add ------------------------------ */

// Adds the specified amount to the current progress and returns the new
// 'current' value.
func (p *Progress) add(n uint64) uint64 {
	return p.current.Add(n)
}
