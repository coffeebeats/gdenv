package extract_test

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/coffeebeats/gdenv/internal/extract"
)

// entry describes a single archive member. Names deliberately use '/' - the
// separator mandated by both the 'tar' and 'zip' formats - rather than
// 'filepath.Join', so that these fixtures behave identically on Windows.
type entry struct {
	name     string
	contents string
	mode     fs.FileMode
	isDir    bool
}

/* -------------------------------------------------------------------------- */
/*                            Function: TestTarGz                             */
/* -------------------------------------------------------------------------- */

func TestTarGzExtract(t *testing.T) {
	tests := []struct {
		name    string
		entries []entry
		opts    []extract.Option
		want    map[string]string
		wantErr error
	}{
		{
			name: "flat files are extracted",
			entries: []entry{
				{name: "gdenv", contents: "gdenv-binary", mode: 0o755},
				{name: "godot", contents: "shim-binary", mode: 0o755},
			},
			want: map[string]string{"gdenv": "gdenv-binary", "godot": "shim-binary"},
		},
		{
			name: "nested files are extracted",
			entries: []entry{
				{name: "inner/", mode: 0o755, isDir: true},
				{name: "inner/godot", contents: "shim-binary", mode: 0o755},
			},
			want: map[string]string{"inner/godot": "shim-binary"},
		},
		{
			name: "prefix is stripped when requested",
			entries: []entry{
				{name: "pkg/", mode: 0o755, isDir: true},
				{name: "pkg/gdenv", contents: "gdenv-binary", mode: 0o755},
			},
			opts: []extract.Option{extract.WithStripPrefix("pkg")},
			want: map[string]string{"gdenv": "gdenv-binary"},
		},
		{
			// Guards the empty-prefix case: 'strings.HasPrefix(name, "")' is
			// always true, so an unguarded strip would reject every entry.
			name:    "empty prefix leaves names unchanged",
			entries: []entry{{name: "gdenv", contents: "gdenv-binary", mode: 0o755}},
			opts:    []extract.Option{extract.WithStripPrefix("")},
			want:    map[string]string{"gdenv": "gdenv-binary"},
		},
		{
			name:    "parent-relative path is rejected",
			entries: []entry{{name: "../escape", contents: "evil", mode: 0o644}},
			wantErr: tar.ErrInsecurePath,
		},
		{
			name:    "backslash-separated path is rejected",
			entries: []entry{{name: `inner\escape`, contents: "evil", mode: 0o644}},
			wantErr: tar.ErrInsecurePath,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmp := t.TempDir()
			archive := filepath.Join(tmp, "archive.tar.gz")
			writeTarGz(t, archive, tc.entries)

			out := filepath.Join(tmp, "out")
			if err := os.Mkdir(out, 0o755); err != nil {
				t.Fatal(err)
			}

			err := extract.TarGz(context.Background(), archive, out, tc.opts...)

			assertExtracted(t, out, err, tc.want, tc.wantErr)
		})
	}
}

/* -------------------------------------------------------------------------- */
/*                             Function: TestZip                              */
/* -------------------------------------------------------------------------- */

func TestZipExtract(t *testing.T) {
	tests := []struct {
		name    string
		entries []entry
		want    map[string]string
		wantErr error
	}{
		{
			name: "flat files are extracted",
			entries: []entry{
				{name: "gdenv.exe", contents: "gdenv-binary", mode: 0o755},
				{name: "godot.exe", contents: "shim-binary", mode: 0o755},
			},
			want: map[string]string{"gdenv.exe": "gdenv-binary", "godot.exe": "shim-binary"},
		},
		{
			name: "nested files are extracted",
			entries: []entry{
				{name: "inner/", mode: 0o755, isDir: true},
				{name: "inner/godot.exe", contents: "shim-binary", mode: 0o755},
			},
			want: map[string]string{"inner/godot.exe": "shim-binary"},
		},
		{
			name:    "parent-relative path is rejected",
			entries: []entry{{name: "../escape", contents: "evil", mode: 0o644}},
			wantErr: zip.ErrInsecurePath,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmp := t.TempDir()
			archive := filepath.Join(tmp, "archive.zip")
			writeZip(t, archive, tc.entries)

			out := filepath.Join(tmp, "out")
			if err := os.Mkdir(out, 0o755); err != nil {
				t.Fatal(err)
			}

			err := extract.Zip(context.Background(), archive, out)

			assertExtracted(t, out, err, tc.want, tc.wantErr)
		})
	}
}

/* -------------------------------------------------------------------------- */
/*                      Function: TestTarGzPreservesMode                      */
/* -------------------------------------------------------------------------- */

// TestTarGzPreservesMode verifies the executable bit survives extraction, which
// is what makes a self-updated binary runnable.
func TestTarGzPreservesMode(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("file modes are not meaningfully preserved on Windows")
	}

	tmp := t.TempDir()
	archive := filepath.Join(tmp, "archive.tar.gz")
	writeTarGz(t, archive, []entry{{name: "gdenv", contents: "binary", mode: 0o755}})

	out := filepath.Join(tmp, "out")
	if err := os.Mkdir(out, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := extract.TarGz(context.Background(), archive, out); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(filepath.Join(out, "gdenv"))
	if err != nil {
		t.Fatal(err)
	}

	if got := info.Mode().Perm(); got != 0o755 {
		t.Fatalf("mode: got %o, want %o", got, 0o755)
	}
}

/* -------------------------------------------------------------------------- */
/*                            Function: assertExtracted                       */
/* -------------------------------------------------------------------------- */

func assertExtracted(t *testing.T, out string, err error, want map[string]string, wantErr error) {
	t.Helper()

	if wantErr != nil {
		if !errors.Is(err, wantErr) {
			t.Fatalf("err: got %v, want %v", err, wantErr)
		}

		return
	}

	if err != nil {
		t.Fatalf("err: got %v, want nil", err)
	}

	for name, contents := range want {
		got, err := os.ReadFile(filepath.Join(out, filepath.FromSlash(name)))
		if err != nil {
			t.Fatalf("reading %s: %v", name, err)
		}

		if string(got) != contents {
			t.Fatalf("%s: got %q, want %q", name, string(got), contents)
		}
	}
}

/* -------------------------------------------------------------------------- */
/*                            Function: writeTarGz                            */
/* -------------------------------------------------------------------------- */

func writeTarGz(t *testing.T, path string, entries []entry) {
	t.Helper()

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}

	defer f.Close()

	gw := gzip.NewWriter(f)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	for _, e := range entries {
		hdr := &tar.Header{ //nolint:exhaustruct_v5
			Name:     e.name,
			Mode:     int64(e.mode),
			Typeflag: tar.TypeReg,
			Size:     int64(len(e.contents)),
		}

		if e.isDir {
			hdr.Typeflag = tar.TypeDir
			hdr.Size = 0
		}

		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatal(err)
		}

		if e.isDir {
			continue
		}

		if _, err := tw.Write([]byte(e.contents)); err != nil {
			t.Fatal(err)
		}
	}
}

/* -------------------------------------------------------------------------- */
/*                             Function: writeZip                             */
/* -------------------------------------------------------------------------- */

func writeZip(t *testing.T, path string, entries []entry) {
	t.Helper()

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}

	defer f.Close()

	zw := zip.NewWriter(f)
	defer zw.Close()

	for _, e := range entries {
		hdr := &zip.FileHeader{Name: e.name} //nolint:exhaustruct_v5
		hdr.SetMode(e.mode)

		w, err := zw.CreateHeader(hdr)
		if err != nil {
			t.Fatal(err)
		}

		if e.isDir {
			continue
		}

		if _, err := w.Write([]byte(e.contents)); err != nil {
			t.Fatal(err)
		}
	}
}
