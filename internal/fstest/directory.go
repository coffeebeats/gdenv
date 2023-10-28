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
	Path string
}

/* ----------------------------- Impl: Asserter ----------------------------- */

func (d Dir) Assert(t *testing.T, tempDir string) {
	t.Helper()

	path := clean(t, tempDir, d.Path)

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

func (d Dir) Write(t *testing.T, tempDir string) {
	t.Helper()

	path := clean(t, tempDir, d.Path)

	if err := os.MkdirAll(path, osutil.ModeUserRWXGroupRX); err != nil {
		t.Fatalf("%s: failed to write directory: %s", err, path)
	}
}
