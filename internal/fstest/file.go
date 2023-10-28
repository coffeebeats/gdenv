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
	Path, Contents string
}

/* ----------------------------- Impl: Asserter ----------------------------- */

func (f File) Assert(t *testing.T, tempDir string) {
	t.Helper()

	path := clean(t, tempDir, f.Path)

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

func (f File) Write(t *testing.T, tempDir string) {
	t.Helper()

	path := clean(t, tempDir, f.Path)

	if err := os.MkdirAll(filepath.Dir(path), osutil.ModeUserRWXGroupRX); err != nil {
		t.Fatalf("%s: failed to write directory: %s", err, path)
	}

	if err := os.WriteFile(path, []byte(f.Contents), osutil.ModeUserRW); err != nil {
		t.Fatalf("%s: failed to write file: %s", err, path)
	}
}
