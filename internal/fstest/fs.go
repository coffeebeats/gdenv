package fstest

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
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

	path := a.Path
	if !filepath.IsAbs(path) {
		path = filepath.Join(tempDir, path)
	}

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
/*                                Struct: File                                */
/* -------------------------------------------------------------------------- */

type File struct {
	Path, Contents string
}

/* ----------------------------- Impl: Asserter ----------------------------- */

func (f File) Assert(t *testing.T, tempDir string) {
	t.Helper()

	path := f.Path
	if !filepath.IsAbs(path) {
		path = filepath.Join(tempDir, path)
	}

	info, err := os.Stat(path)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			t.Fatal(err)
		}

		t.Fatalf("file not found: %s", path)
	}

	if !info.Mode().IsRegular() {
		t.Fatalf("expected a file, got: %v", info.Mode().Type())
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

	path := f.Path
	if !filepath.IsAbs(path) {
		path = filepath.Join(tempDir, path)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		t.Fatalf("%s: failed to write directory: %s", err, path)
	}

	if err := os.WriteFile(path, []byte(f.Contents), 0600); err != nil {
		t.Fatalf("%s: failed to write file: %s", err, path)
	}
}

/* -------------------------------------------------------------------------- */
/*                                 Struct: Dir                                */
/* -------------------------------------------------------------------------- */

type Dir struct {
	Path string
}

/* ----------------------------- Impl: Asserter ----------------------------- */

func (d Dir) Assert(t *testing.T, tempDir string) {
	t.Helper()

	path := d.Path
	if !filepath.IsAbs(path) {
		path = filepath.Join(tempDir, path)
	}

	info, err := os.Stat(path)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			t.Fatal(err)
		}

		t.Fatalf("directory not found: %s", path)
	}

	if !info.IsDir() {
		t.Errorf("expected a directory, got: %v", info.Mode().Type())
	}
}

/* ------------------------------ Impl: Writer ------------------------------ */

func (d Dir) Write(t *testing.T, tempDir string) {
	t.Helper()

	path := d.Path
	if !filepath.IsAbs(path) {
		path = filepath.Join(tempDir, path)
	}

	if err := os.MkdirAll(path, 0700); err != nil {
		t.Fatalf("%s: failed to write directory: %s", err, path)
	}
}
