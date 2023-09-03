package progress

import (
	"errors"
	"io"
)

var ErrMissingProgress = errors.New("missing progress")

/* -------------------------------------------------------------------------- */
/*                               Struct: Writer                               */
/* -------------------------------------------------------------------------- */

// A thread-safe 'io.Writer' implementation that tracks the percentage of bytes
// written against a specified total.
type Writer struct {
	progress *Progress
}

// Validate at compile-time that 'Writer' implements 'io.Writer'.
var _ io.Writer = &Writer{} //nolint:exhaustruct

/* --------------------------- Function: NewWriter -------------------------- */

// Creates a new 'Writer' with the specified 'Progress' reporter.
//
// NOTE: It's the caller's responsibility to ensure that the initial 'total'
// size is correct so that the computed progress value is accurate.
func NewWriter(p *Progress) Writer {
	return Writer{p}
}

/* ----------------------------- Impl: io.Writer ---------------------------- */

func (w Writer) Write(data []byte) (int, error) {
	// NOTE: This cannot be silently corrected (i.e. create 'Progress' here)
	// because 'Writer' may have been copied prior to this.
	if w.progress == nil {
		return 0, ErrMissingProgress
	}

	n := len(data)

	if _, err := w.progress.add(uint64(n)); err != nil {
		return 0, err
	}

	return n, nil
}
