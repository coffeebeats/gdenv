package progress

import "io"

/* -------------------------------------------------------------------------- */
/*                               Struct: Writer                               */
/* -------------------------------------------------------------------------- */

// A thread-safe 'io.Writer' implementation that tracks the percentage of bytes
// written against a specified total.
type Writer struct {
	progress *Progress
}

// Validate at compile-time that 'Writer' implements 'io.Writer'.
var _ io.Writer = (*Writer)(nil)

/* --------------------------- Function: NewWriter -------------------------- */

// Creates a new 'Writer' with the specified 'Progress' reporter.
//
// NOTE: It's the caller's responsibility to ensure that the initial 'Progress'
// provided is correctly configured so that the calculated progress is accurate.
func NewWriter(p *Progress) *Writer {
	return &Writer{p}
}

/* ----------------------------- Impl: io.Writer ---------------------------- */

func (w *Writer) Write(data []byte) (int, error) {
	n := len(data)

	w.progress.Add(uint64(n))

	return n, nil
}
