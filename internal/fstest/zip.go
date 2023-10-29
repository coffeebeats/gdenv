package fstest

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/coffeebeats/gdenv/internal/osutil"
)

/* -------------------------------------------------------------------------- */
/*                                 Struct: Zip                                */
/* -------------------------------------------------------------------------- */

type Zip struct {
	Path     string
	Contents []Writer
}

/* ------------------------------ Impl: Writer ------------------------------ */

func (z Zip) Abs(t *testing.T, pathBaseDir string) string {
	t.Helper()

	return clean(t, pathBaseDir, z.Path)
}

func (z Zip) Write(t *testing.T, pathBaseDir string) {
	t.Helper()

	// First, create the archive in its own temporary directory.
	tmpArchive := z.Abs(t, t.TempDir())
	createArchiveWithContents(t, z.Contents, tmpArchive)

	// Then, copy the archive and update the permissions as required.
	copyArchiveAndSetPermissions(t, tmpArchive, z.Abs(t, pathBaseDir))
}

/* ----------------- Function: copyArchiveAndSetPermissions ----------------- */

func copyArchiveAndSetPermissions(t *testing.T, path, out string) {
	t.Helper()

	r, err := zip.OpenReader(path)
	if err != nil {
		t.Fatal(err)
	}

	defer r.Close()

	f, err := os.OpenFile(out, os.O_WRONLY|os.O_CREATE, osutil.ModeUserRWX)
	if err != nil {
		t.Fatalf("test setup: failed to create file: %v", path)
	}

	defer f.Close()

	w := zip.NewWriter(f)

	defer w.Close()

	for _, f := range r.File {
		if f.Mode().IsDir() {
			f.SetMode(osutil.ModeUserRWX)
		} else {
			f.SetMode(osutil.ModeUserRW)
		}

		if err := w.Copy(f); err != nil {
			t.Fatal(err)
		}
	}
}

/* ------------------- Function: createArchiveWithContents ------------------ */

func createArchiveWithContents(t *testing.T, contents []Writer, out string) {
	t.Helper()

	f, err := os.OpenFile(out, os.O_WRONLY|os.O_CREATE, osutil.ModeUserRWX)
	if err != nil {
		t.Fatalf("test setup: failed to create file: %v", out)
	}

	defer f.Close()

	a := zip.NewWriter(f)

	defer a.Close()

	for _, c := range contents {
		switch w := c.(type) {
		case Dir:
			createZipDirectory(t, a, out, clean(t, out, w.Path))
		case File:
			createZipFile(t, a, out, clean(t, out, w.Path), w.Contents)
		default:
			t.Fatalf("unsupported content type: %T", w)
		}
	}
}

/* ---------------------- Function: createZipDirectory ---------------------- */

func createZipDirectory(t *testing.T, archive *zip.Writer, pathBaseDir, name string) {
	t.Helper()

	p, err := filepath.Rel(pathBaseDir, name)
	if err != nil {
		t.Fatalf("test setup: failed to determine path: %v", err)
	}

	if _, err := archive.Create(p + string(os.PathSeparator)); err != nil {
		t.Fatalf("test setup: failed to create directory: %v", err)
	}
}

/* ------------------------- Function: createZipFile ------------------------ */

func createZipFile(t *testing.T, archive *zip.Writer, pathBaseDir, name, contents string) {
	t.Helper()

	p, err := filepath.Rel(pathBaseDir, name)
	if err != nil {
		t.Fatalf("test setup: failed to determine path: %v", err)
	}

	dst, err := archive.Create(p)
	if err != nil {
		t.Fatalf("test setup: failed to create file: %v", err)
	}

	if _, err := io.Copy(dst, strings.NewReader(contents)); err != nil {
		t.Fatalf("test setup: failed to write to file: %v", err)
	}
}
