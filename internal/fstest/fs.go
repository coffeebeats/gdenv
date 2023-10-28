package fstest

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

/* -------------------------------------------------------------------------- */
/*                             Interface: Asserter                            */
/* -------------------------------------------------------------------------- */

type Asserter interface {
	Assert(t *testing.T, tempDir string)
}

/* -------------------------------------------------------------------------- */
/*                              Interface: Writer                             */
/* -------------------------------------------------------------------------- */

type Writer interface {
	Write(t *testing.T, tempDir string)
}

/* -------------------------------------------------------------------------- */
/*                               Struct: Absent                               */
/* -------------------------------------------------------------------------- */

type Absent struct {
	Path string
}

/* ----------------------------- Impl: Asserter ----------------------------- */

func (a Absent) Assert(t *testing.T, tempDir string) {
	t.Helper()

	path := clean(t, tempDir, a.Path)

	info, err := os.Stat(path)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			t.Fatal(err)
		}

		return
	}

	t.Errorf("unexpectedly found file: %s (%v)", path, info.Mode().Type())
}

/* ------------------------------ Impl: Writer ------------------------------ */

func (a Absent) Write(t *testing.T, _ string) {
	t.Helper()
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

	info, err := os.Stat(base)
	if err != nil {
		t.Fatal(err)
	}

	if !info.IsDir() {
		t.Fatalf("expected 'base' path to be a directory: %s", base)
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
