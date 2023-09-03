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

// A thread-safe progress tracker; does minimal error-case handling, will wrap
// any invalid values around '0', and reports percentages exactly (i.e. greater
// than '1.0' is possible).
//
// NOTE: Avoid copying this struct, as it will be unlinked from other consumers.
type Progress struct {
	current atomic.Uint64
	total   float64
}

/* ------------------------------ Function: New ----------------------------- */

// Creates a new 'Progress' struct with the 'total' bounded above '0'.
func New(total int64) (Progress, error) {
	if total <= 0 {
		return Progress{}, fmt.Errorf("%w: %d", ErrInvalidTotal, total)
	}

	return Progress{current: atomic.Uint64{}, total: float64(total)}, nil
}

/* --------------------------- Method: Percentage --------------------------- */

// Retrieves the current progress as a percentage in the range [0.0, 1.0].
func (p *Progress) Percentage() float64 {
	if p.total <= 0 {
		return 0
	}

	return float64(p.current.Load()) / p.total
}

/* ------------------------------- Method: add ------------------------------ */

// Adds the specified amount to the current progress.
func (p *Progress) add(n uint64) uint64 {
	return p.current.Add(n)
}
