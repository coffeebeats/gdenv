package progress

import "io"

/* -------------------------------------------------------------------------- */
/*                            Struct: ManualWriter                            */
/* -------------------------------------------------------------------------- */

// A thread-safe 'io.Writer' implementation that allows for manual tracking of
// a percentage of progress made against a total.
type ManualWriter struct {
	progress *Progress
}

/* --------------------------- Function: NewWriter -------------------------- */

// Creates a new 'ManualWriter' with the specified 'Progress' reporter.
//
// NOTE: It's the caller's responsibility to ensure that the initial 'Progress'
// provided is correctly configured so that the calculated progress is accurate.
func NewManualWriter(p *Progress) *ManualWriter {
	return &ManualWriter{p}
}

/* ------------------------------- Method: Add ------------------------------ */

// Adds the specified amount to the current progress and returns the new
// 'current' value.
func (w *ManualWriter) Add(n uint64) uint64 {
	return w.progress.add(n)
}

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
// NOTE: It's the caller's responsibility to ensure that the initial 'Progress'
// provided is correctly configured so that the calculated progress is accurate.
func NewWriter(p *Progress) *Writer {
	return &Writer{p}
}

/* ----------------------------- Impl: io.Writer ---------------------------- */

func (w Writer) Write(data []byte) (int, error) {
	n := len(data)

	w.progress.add(uint64(n))

	return n, nil
}
