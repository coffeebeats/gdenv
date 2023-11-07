package fstest

import (
	"path/filepath"
	"testing"
)

/* -------------------------------------------------------------------------- */
/*                             Interface: Filepath                            */
/* -------------------------------------------------------------------------- */

type Filepath interface {
	Resolve(t *testing.T, pathBaseDir string) string
}

/* -------------------------------------------------------------------------- */
/*                                 Type: Exact                                */
/* -------------------------------------------------------------------------- */

// Exact is a wrapper around a filepath that is returned as-is.
type Exact string

// Compile-time check that 'Filepath' is implemented.
var _ Filepath = (*Exact)(nil)

/* -------------------------- Impl: fstest.Filepath ------------------------- */

func (e Exact) Resolve(t *testing.T, _ string) string {
	t.Helper()

	return string(e)
}

/* -------------------------------------------------------------------------- */
/*                               Type: Relative                               */
/* -------------------------------------------------------------------------- */

// Relative is a wrapper around a filepath that returns a relative path.
type Relative string

// Compile-time check that 'Filepath' is implemented.
var _ Filepath = (*Relative)(nil)

/* -------------------------- Impl: fstest.Filepath ------------------------- */

// Filepath returns a filepath relative to the provided base directory. Fails if
// the underlying path cannot be made relative to the base directory.
func (r Relative) Resolve(t *testing.T, pathBaseDir string) string {
	t.Helper()

	path := filepath.Clean(string(r))
	if !filepath.IsAbs(path) {
		return path
	}

	path, err := filepath.Rel(pathBaseDir, path)
	if err != nil {
		t.Fatal(err)
	}

	return path
}

/* -------------------------------------------------------------------------- */
/*                               Type: Absolute                               */
/* -------------------------------------------------------------------------- */

// Absolute is a wrapper around a filepath that returns an absolute path.
type Absolute string

// Compile-time check that 'Filepath' is implemented.
var _ Filepath = (*Absolute)(nil)

/* -------------------------- Impl: fstest.Filepath ------------------------- */

// Filepath returns an absolute path to the wrapped path. If the underlying path
// is absolute then it must be a descendent of the provided base directory.
func (a Absolute) Resolve(t *testing.T, pathBaseDir string) string {
	t.Helper()

	path := filepath.Clean(string(a))
	if !filepath.IsAbs(path) {
		path = filepath.Join(pathBaseDir, path)
	}

	if _, err := filepath.Rel(pathBaseDir, path); err != nil {
		t.Fatal(err)
	}

	return path
}
