package progress

import (
	"math"
	"sync/atomic"
)

/* -------------------------------------------------------------------------- */
/*                              Struct: Progress                              */
/* -------------------------------------------------------------------------- */

// A thread-safe progress tracker; does minimal error-case handling, will wrap
// any invalid values around '0', and reports percentages exactly (i.e. greater
// than '1.0' is possible).
type Progress struct {
	current atomic.Uint64
	total   float64
}

/* -------------------------- Function: NewProgress ------------------------- */

// Creates a new 'Progress' struct with the 'size' bounded above '0'.
func NewProgress(size int64) Progress {
	return Progress{total: math.Max(float64(size), 0)} //nolint:exhaustruct
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
func (p *Progress) add(n uint64) {
	p.current.Add(n)
}
