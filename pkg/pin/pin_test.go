package pin

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/coffeebeats/gdenv/internal/fstest"
	"github.com/coffeebeats/gdenv/internal/osutil"
	"github.com/coffeebeats/gdenv/pkg/godot/version"
)

/* ------------------------------- Test: Read ------------------------------- */

func TestRead(t *testing.T) {
	tests := []struct {
		name string

		path  fstest.Filepath
		files []fstest.Writer

		want version.Version
		err  error
	}{
		// Invalid inputs
		{
			name: "missing path returns an error",
			path: fstest.Exact(""),
			err:  ErrMissingPath,
		},
		{
			name:  "invalid pin file returns an error",
			path:  fstest.Absolute(".godot-version"),
			files: []fstest.Writer{fstest.File{Path: ".godot-version", Contents: "invalid"}},
			err:   version.ErrInvalid,
		},
		{
			name: "missing file returns an error",
			path: fstest.Relative(".godot-version"),
			err:  ErrMissingPin,
		},

		// Valid inputs
		{
			name:  "specifying an exact file reads it",
			path:  fstest.Absolute(".godot-version"),
			files: []fstest.Writer{fstest.File{Path: ".godot-version", Contents: "v1.0-stable"}},
			want:  version.MustParse("1.0.0-stable"),
		},
		{
			name:  "specifying a directory reads the file within that directory",
			path:  fstest.Absolute("a/b/c"),
			files: []fstest.Writer{fstest.File{Path: "a/b/c/.godot-version", Contents: "v1.0-stable"}},
			want:  version.MustParse("1.0.0-stable"),
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			tmp := t.TempDir()

			// Given: The specified files exist on the file system.
			for _, f := range tc.files {
				f.Write(t, tmp)
			}

			got, err := Read(tc.path.Resolve(t, tmp))
			if !errors.Is(err, tc.err) {
				t.Errorf("err: got %v, want %v", err, tc.err)
			}

			if got != tc.want {
				t.Errorf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}

}

/* ----------------------------- Test: VersionAt ---------------------------- */

func TestVersionAt(t *testing.T) {
	tests := []struct {
		pin  string // where the pin file exists
		path string // where to query

		want version.Version // result of the query
		err  error
	}{
		{pin: "a", path: "", err: ErrMissingPin},
		{pin: "a/b", path: "c/d", err: ErrMissingPin},

		{pin: "", path: "", want: version.Godot3()},
		{pin: "", path: "a/b/c", want: version.Godot3()},
		{pin: ".gdenv", path: "a/b/c", want: version.Godot3()}, // uses global pin
	}

	for i, tc := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			tmp := t.TempDir()

			path, pin := filepath.Join(tmp, tc.path), filepath.Join(tmp, tc.pin)

			pin, err := clean(pin)
			if err != nil {
				t.Fatalf("test setup: %v", err)
			}

			// Create the pin file
			if err := os.MkdirAll(filepath.Dir(pin), osutil.ModeUserRWXGroupRX); err != nil {
				t.Fatalf("test setup: %v", err)
			}

			if err := os.WriteFile(pin, []byte(tc.want.String()), osutil.ModeUserRW); err != nil {
				t.Fatalf("test setup: %v", err)
			}

			storePath := filepath.Join(tmp, ".gdenv")
			got, err := VersionAt(context.Background(), storePath, path)
			if !errors.Is(err, tc.err) {
				t.Errorf("err: got %v, want %v", err, tc.err)
			}

			if got != tc.want {
				t.Errorf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}

/* ------------------------------ Test: Remove ------------------------------ */

func TestRemove(t *testing.T) {
	tests := []struct {
		path     string
		existing bool
		err      error
	}{
		{"", false, nil},
		{"", true, nil},
		{"a/b/c", true, nil},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			tmp := t.TempDir()
			path := filepath.Join(tmp, tc.path)

			pin, err := clean(path)
			if err != nil {
				t.Fatalf("test setup: %v", err)
			}

			if tc.existing {
				if err := os.MkdirAll(filepath.Dir(pin), osutil.ModeUserRWXGroupRX); err != nil {
					t.Fatalf("test setup: %v", err)
				}

				if err := os.WriteFile(pin, []byte(""), osutil.ModeUserRW); err != nil {
					t.Fatalf("test setup: %v", err)
				}
			}

			err = Remove(path)
			if !errors.Is(err, tc.err) {
				t.Errorf("err: got %v, want %v", err, tc.err)
			}

			if _, err := os.Stat(pin); !errors.Is(err, fs.ErrNotExist) {
				t.Errorf("err: %v", err)
			}
		})
	}

}

/* ------------------------------- Test: Write ------------------------------ */

func TestWrite(t *testing.T) {
	tests := []struct {
		version string
		path    func(tempDir string) string

		want fstest.Asserter // will have 'tempDir' prefixed.
		err  error
	}{
		{
			version: "v4",
			path:    func(_ string) string { return "" },

			err: ErrMissingPath,
		},
		{
			version: "v4",
			path: func(tmp string) string {
				return filepath.Join(tmp, "does/not/exist", pinFilename)
			},

			err: fs.ErrNotExist,
		},
		{
			version: "v4",
			path: func(tmp string) string {
				return filepath.Join(tmp, pinFilename)
			},

			want: fstest.File{Path: pinFilename, Contents: "v4.0-stable"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.version, func(t *testing.T) {
			tmp := t.TempDir()

			// When: The pin file is written are cleared from the store.
			// Then: The expected error value is returned.
			if err := Write(version.MustParse(tc.version), tc.path(tmp)); !errors.Is(err, tc.err) {
				t.Errorf("got: %v, want: %v", err, tc.err)
			}

			// Then: The expected files exist on the file system.
			if tc.want != nil {
				tc.want.Assert(t, tmp)
			}
		})
	}
}
