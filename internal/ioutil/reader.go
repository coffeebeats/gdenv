package ioutil

import "context"

/* -------------------------------------------------------------------------- */
/*                             Type: ReaderClosure                            */
/* -------------------------------------------------------------------------- */

// ReaderClosure wraps a function which implements 'io.Reader', allowing for
// inline 'io.Reader' definitions.
type ReaderClosure func([]byte) (int, error)

/* ------------------------ Function: newReaderCloser ----------------------- */

// NewReaderClosure returns a new 'io.Reader' implementation which cancels
// reading once the provided 'context.Context' is closed.
func NewReaderClosure(ctx context.Context, r func([]byte) (int, error)) ReaderClosure {
	return ReaderClosure(func(p []byte) (int, error) {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
			return r(p)
		}
	})
}

/* ----------------------------- Impl: io.Reader ---------------------------- */

func (r ReaderClosure) Read(p []byte) (int, error) { return r(p) }
