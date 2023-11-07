package fstest

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/coffeebeats/gdenv/internal/osutil"
)

/* -------------------------------------------------------------------------- */
/*                                Struct: File                                */
/* -------------------------------------------------------------------------- */

type File struct {
	// Path is a filepath that will be treated relative to a base directory.
	Path string
	// Contents are written to the created file.
	Contents string
}

/* ----------------------------- Impl: Asserter ----------------------------- */

func (f File) Assert(t *testing.T, pathBaseDir string) {
	t.Helper()

	path := f.Abs(t, pathBaseDir)

	info, err := os.Stat(path)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			t.Fatal(err)
		}

		t.Errorf("file not found: %s", path)

		return
	}

	if !info.Mode().IsRegular() {
		t.Errorf("expected a file, got: %v", info.Mode().Type())

		return
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := string(got), f.Contents; got != want {
		t.Errorf("file contents: got %s, want %s", got, want)
	}
}

/* ------------------------------ Impl: Writer ------------------------------ */

func (f File) Abs(t *testing.T, pathBaseDir string) string {
	t.Helper()

	return clean(t, pathBaseDir, f.Path)
}

func (f File) Write(t *testing.T, pathBaseDir string) {
	t.Helper()

	path := f.Abs(t, pathBaseDir)

	if err := os.MkdirAll(filepath.Dir(path), osutil.ModeUserRWXGroupRX); err != nil {
		t.Fatalf("%s: failed to write directory: %s", err, path)
	}

	if err := os.WriteFile(path, []byte(f.Contents), osutil.ModeUserRW); err != nil {
		t.Fatalf("%s: failed to write file: %s", err, path)
	}
}
