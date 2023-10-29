package ioutil

import (
	"context"
	"io"
)

/* -------------------------------------------------------------------------- */
/*                             Type: ReaderClosure                            */
/* -------------------------------------------------------------------------- */

// ReaderClosure wraps a function which implements 'io.Reader', allowing for
// inline 'io.Reader' definitions.
type ReaderClosure func([]byte) (int, error)

/* ----------------------------- Impl: io.Reader ---------------------------- */

func (r ReaderClosure) Read(p []byte) (int, error) { return r(p) }

/* -------------------------------------------------------------------------- */
/*                       Function: NewReaderWithContext                       */
/* -------------------------------------------------------------------------- */

// NewReaderWithContext returns a new 'io.Reader' implementation which cancels
// reading once the provided 'context.Context' is closed.
func NewReaderWithContext(ctx context.Context, r ReaderClosure) io.Reader {
	return ReaderClosure(func(p []byte) (int, error) {
		if ctx.Err() != nil {
			return 0, ctx.Err()
		}

		return r(p)
	})
}
