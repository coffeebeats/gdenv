package ioutil

import (
	"bytes"
	"io"
)

/* -------------------------------------------------------------------------- */
/*                             Type: WriterClosure                            */
/* -------------------------------------------------------------------------- */

// WriterClosure wraps a function which implements 'io.Writer', allowing for
// inline 'io.Writer' definitions.
type WriterClosure func([]byte) (int, error)

/* ----------------------------- Impl: io.Writer ---------------------------- */

func (w WriterClosure) Write(p []byte) (int, error) { return w(p) }

/* -------------------------------------------------------------------------- */
/*                         Function: NewWriterWithTrim                        */
/* -------------------------------------------------------------------------- */

// NewWriterWithTrim returns a new 'io.Writer' implementation which trims null
// bytes from the input byte slice.
func NewWriterWithTrim(w io.Writer, cutset string) io.Writer {
	return WriterClosure(func(p []byte) (int, error) {
		t := bytes.Trim(p, cutset)

		_, err := w.Write(t)
		if err != nil {
			return 0, err
		}

		// The 'io.Writer' interface must return a non-nil error if the returned
		// value of 'n' is less than 'len(p)'. As such, return the full length
		// of 'p' despite writing fewer bytes to the wrapped 'io.Writer'.
		return len(p), nil
	})
}
