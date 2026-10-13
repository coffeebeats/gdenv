package fstest

import (
	"archive/tar"
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/ulikunitz/xz"

	"github.com/coffeebeats/gdenv/internal/osutil"
)

/* -------------------------------------------------------------------------- */
/*                                Struct: TarXZ                               */
/* -------------------------------------------------------------------------- */

type TarXZ struct {
	// Path is a filepath that will be treated relative to a base directory.
	Path     string
	Contents []Writer
}

/* ------------------------------ Impl: Writer ------------------------------ */

func (tr TarXZ) Abs(t *testing.T, pathBaseDir string) string {
	t.Helper()

	return clean(t, pathBaseDir, tr.Path)
}

func (tr TarXZ) Write(t *testing.T, pathBaseDir string) {
	t.Helper()

	path := tr.Abs(t, pathBaseDir)

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, osutil.ModeUserRWX)
	if err != nil {
		t.Fatalf("test setup: failed to create file: %v", path)
	}

	defer f.Close()

	w, err := xz.NewWriter(f)
	if err != nil {
		t.Fatalf("test setup: failed to create xz writer: %v", err)
	}

	defer w.Close()

	tw := tar.NewWriter(w)

	defer tw.Close()

	for _, c := range tr.Contents {
		var contents []byte

		hdr := &tar.Header{} //nolint:exhaustruct_v5

		switch c := c.(type) {
		case Dir:
			hdr.Name = entry(t, path, c.Path)
			hdr.Mode = osutil.ModeUserRWX
			hdr.Typeflag = tar.TypeDir
		case File:
			contents = []byte(c.Contents)

			hdr.Name = entry(t, path, c.Path)
			hdr.Mode = osutil.ModeUserRW
			hdr.Typeflag = tar.TypeReg
			hdr.Size = int64(len(contents))

		default:
			t.Fatalf("unsupported content type: %T", w)
		}

		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatal(err)
		}

		if hdr.Typeflag != tar.TypeReg {
			continue
		}

		if _, err := io.Copy(tw, bytes.NewReader(contents)); err != nil {
			t.Fatal(err)
		}
	}
}
