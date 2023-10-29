package fstest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

/* -------------------------------------------------------------------------- */
/*                             Interface: Asserter                            */
/* -------------------------------------------------------------------------- */

type Asserter interface {
	Assert(t *testing.T, pathBaseDir string)
}

/* -------------------------------------------------------------------------- */
/*                              Interface: Writer                             */
/* -------------------------------------------------------------------------- */

type Writer interface {
	Abs(t *testing.T, pathBaseDir string) string
	Write(t *testing.T, pathBaseDir string)
}

/* -------------------------------------------------------------------------- */
/*                               Function: clean                              */
/* -------------------------------------------------------------------------- */

// clean is a wrapper around 'filepath.Clean' that ensures (1) the provided path
// is underneath the specified 'base' directory and (2) the returned path is an
// absolute, cleaned path.
func clean(t *testing.T, base, path string) string {
	t.Helper()

	if !filepath.IsAbs(base) {
		t.Fatalf("expected 'base' path to be absolute: %s", base)
	}

	if _, err := os.Stat(base); err != nil {
		t.Fatal(err)
	}

	path = filepath.Clean(path)
	if filepath.IsAbs(path) {
		prefix := base + string(os.PathSeparator)
		if !strings.HasPrefix(path, prefix) {
			t.Fatalf("invalid absolute path: %s", path)
		}

		// NOTE: Because 'path' is first cleaned, a valid 'path' will *not* have
		// a trailing path separator. As such, we can safely trim it off.
		path = strings.TrimPrefix(path, prefix)
	}

	return filepath.Join(base, path)
}
