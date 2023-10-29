package fstest

import (
	"errors"
	"io/fs"
	"os"
	"testing"
)

/* -------------------------------------------------------------------------- */
/*                               Struct: Absent                               */
/* -------------------------------------------------------------------------- */

type Absent struct {
	Path string
}

/* ----------------------------- Impl: Asserter ----------------------------- */

func (a Absent) Assert(t *testing.T, pathBaseDir string) {
	t.Helper()

	path := a.Abs(t, pathBaseDir)

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

func (a Absent) Abs(t *testing.T, pathBaseDir string) string {
	t.Helper()

	return clean(t, pathBaseDir, a.Path)
}

func (a Absent) Write(t *testing.T, _ string) {
	t.Helper()
}
