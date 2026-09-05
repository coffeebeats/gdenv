// Package extract implements archive extraction for the formats used by both
// Godot release artifacts and 'gdenv' release artifacts.
//
// The logic here is format-level and deliberately free of any artifact typing;
// callers which need a typed API (see 'pkg/godot/artifact/archive') wrap these
// functions rather than reimplementing them.
package extract

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"os"

	"github.com/coffeebeats/gdenv/internal/ioutil"
	"github.com/coffeebeats/gdenv/internal/osutil"
	"github.com/coffeebeats/gdenv/pkg/progress"
)

// Only write to 'out'; create a new file/overwrite an existing.
const copyFileWriteFlag = os.O_WRONLY | os.O_CREATE | os.O_TRUNC

var ErrFailed = errors.New("extract failed")

type progressKey struct{}

/* -------------------------------------------------------------------------- */
/*                           Function: WithProgress                           */
/* -------------------------------------------------------------------------- */

// WithProgress creates a sub-context with an associated progress reporter. The
// result can be passed to the extract function(s) in this package to get
// updates on extraction progress.
func WithProgress(ctx context.Context, p *progress.Progress) context.Context {
	return context.WithValue(ctx, progressKey{}, p)
}

/* -------------------------------------------------------------------------- */
/*                               Type: Option                                 */
/* -------------------------------------------------------------------------- */

// An Option configures the behavior of an extract function in this package.
type Option func(*config)

// config collects the configurable behavior of an extract function.
type config struct {
	stripPrefix string
}

/* -------------------------- Function: WithStripPrefix --------------------- */

// WithStripPrefix removes the provided directory prefix from each archived
// path, allowing contents nested under a top-level directory to be extracted
// directly into the output directory.
//
// NOTE: An empty prefix is ignored; archived paths are used as-is.
func WithStripPrefix(prefix string) Option {
	return func(c *config) {
		c.stripPrefix = prefix
	}
}

/* --------------------------- Function: newConfig -------------------------- */

// newConfig applies the provided options to a default configuration.
func newConfig(opts ...Option) config {
	var c config

	for _, opt := range opts {
		opt(&c)
	}

	return c
}

/* --------------------------- Function: copyFile --------------------------- */

// A shared helper function which copies the contents of an 'io.Reader' to a new
// file created with the specified 'os.FileMode'.
func copyFile(ctx context.Context, f io.Reader, mode fs.FileMode, out string) error {
	dst, err := os.OpenFile(out, copyFileWriteFlag, mode)
	if err != nil {
		return err
	}

	defer dst.Close()

	if _, err := io.Copy(dst, ioutil.NewReaderWithContext(ctx, f.Read)); err != nil {
		return err
	}

	return nil
}

/* ------------------- Function: newFileReaderWithProgress ------------------ */

// newFileReaderWithProgress sets the 'total' value of the 'progress.Progress'
// instance attached to the context, if one exists. A pointer to the provided
// 'progress.Progress' is returned.
func newFileReaderWithProgress(ctx context.Context, f *os.File) (io.Reader, error) {
	p, ok := ctx.Value(progressKey{}).(*progress.Progress)
	if !ok || p == nil {
		return f, nil
	}

	sum, err := osutil.SizeOf(f.Name())
	if err != nil {
		return f, err
	}

	if err := p.SetTotal(sum); err != nil {
		return f, err
	}

	return io.TeeReader(f, progress.NewWriter(p)), nil
}
