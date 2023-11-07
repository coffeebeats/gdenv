package fstest

import (
	"errors"
	"io/fs"
	"os"
	"testing"

	"github.com/coffeebeats/gdenv/internal/osutil"
)

/* -------------------------------------------------------------------------- */
/*                                 Struct: Dir                                */
/* -------------------------------------------------------------------------- */

type Dir struct {
	// Path is a filepath that will be treated relative to a base directory.
	Path string
}

/* ----------------------------- Impl: Asserter ----------------------------- */

func (d Dir) Assert(t *testing.T, pathBaseDir string) {
	t.Helper()

	path := d.Abs(t, pathBaseDir)

	info, err := os.Stat(path)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			t.Fatal(err)
		}

		t.Errorf("directory not found: %s", path)

		return
	}

	if !info.IsDir() {
		t.Errorf("expected a directory, got: %v", info.Mode().Type())
	}
}

/* ------------------------------ Impl: Writer ------------------------------ */

func (d Dir) Abs(t *testing.T, pathBaseDir string) string {
	t.Helper()

	return clean(t, pathBaseDir, d.Path)
}

func (d Dir) Write(t *testing.T, pathBaseDir string) {
	t.Helper()

	path := d.Abs(t, pathBaseDir)

	if err := os.MkdirAll(path, osutil.ModeUserRWXGroupRX); err != nil {
		t.Fatalf("%s: failed to write directory: %s", err, path)
	}
}
