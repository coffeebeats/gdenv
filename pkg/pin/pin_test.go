package pin

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/coffeebeats/gdenv/internal/godot"
)

/* ------------------------------- Test: Read ------------------------------- */

func TestRead(t *testing.T) {
	version := godot.Version{}.Canonical()

	tests := []struct {
		path     string
		existing bool
		want     godot.Version
		err      error
	}{
		{"", true, version, nil},
		{"a/b/c", true, version, nil},
		{"", false, godot.Version{}, ErrFileNotFound},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			tmp := t.TempDir()
			path := filepath.Join(tmp, tc.path)

			pin, err := Clean(path)
			if err != nil {
				t.Fatalf("test setup: %v", err)
			}

			if tc.existing {
				// Create the pin file
				if err := os.MkdirAll(filepath.Dir(pin), os.ModePerm); err != nil {
					t.Fatalf("test setup: %v", err)
				}

				if err := os.WriteFile(pin, []byte(version.String()), os.ModePerm); err != nil {
					t.Fatalf("test setup: %v", err)
				}
			}

			got, err := Read(path)
			if !errors.Is(err, tc.err) {
				t.Fatalf("err: got %v, want %v", err, tc.err)
			}

			if got != tc.want {
				t.Fatalf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}

}

/* ------------------------------ Test: Resolve ----------------------------- */

func TestResolve(t *testing.T) {
	tests := []struct {
		pin  string // where the pin file exists
		path string // where to query
		want string // result of the query
		err  error
	}{
		{"", "", pinFilename, nil},
		{"", "a/b/c", pinFilename, nil},
		{"a", "", "", ErrFileNotFound},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			tmp := t.TempDir()

			path, pin := filepath.Join(tmp, tc.path), filepath.Join(tmp, tc.pin)

			pin, err := Clean(pin)
			if err != nil {
				t.Fatalf("test setup: %v", err)
			}

			// Create the pin file
			if err := os.MkdirAll(filepath.Dir(pin), os.ModePerm); err != nil {
				t.Fatalf("test setup: %v", err)
			}

			if err := os.WriteFile(pin, []byte(""), os.ModePerm); err != nil {
				t.Fatalf("test setup: %v", err)
			}

			got, err := Resolve(path)
			if !errors.Is(err, tc.err) {
				t.Fatalf("err: got %v, want %v", err, tc.err)
			}

			got, err = filepath.Rel(tmp, got)
			if got != "" && err != nil {
				t.Fatalf("test setup: %v", err)
			}

			if got != tc.want {
				t.Fatalf("output: got %#v, want %#v", got, tc.want)
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

			pin, err := Clean(path)
			if err != nil {
				t.Fatalf("test setup: %v", err)
			}

			if tc.existing {
				if err := os.MkdirAll(filepath.Dir(pin), os.ModePerm); err != nil {
					t.Fatalf("test setup: %v", err)
				}

				if err := os.WriteFile(pin, []byte(""), os.ModePerm); err != nil {
					t.Fatalf("test setup: %v", err)
				}
			}

			err = Remove(path)
			if !errors.Is(err, tc.err) {
				t.Fatalf("err: got %v, want %v", err, tc.err)
			}

			if _, err := os.Stat(pin); !errors.Is(err, fs.ErrNotExist) {
				t.Fatalf("err: %v", err)
			}
		})
	}

}

/* ------------------------------- Test: Write ------------------------------ */

func TestWrite(t *testing.T) {
	tests := []struct {
		version string
		path    string
		want    string
		err     error
	}{
		{"v4", "", "v4.0.0-stable", nil},
		{"v4", pinFilename, "v4.0.0-stable", nil},
		{"v4.1-rc1", "a/b/c", "v4.1.0-rc1", nil},
		{"v4.1-rc1", "a/b/c/" + pinFilename, "v4.1.0-rc1", nil},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			tmp := t.TempDir()

			v, err := godot.ParseVersion(tc.version)
			if err != nil {
				t.Fatalf("test setup: %v", err)
			}

			p := filepath.Join(tmp, tc.path)

			err = Write(v, p)
			if !errors.Is(err, tc.err) {
				t.Fatalf("err: got %v, want %v", err, tc.err)
			}

			p, err = Clean(p)
			if err != nil {
				t.Fatalf("test setup: %v", err)
			}

			contents, err := os.ReadFile(p)
			if err != nil {
				t.Fatalf("test: %v", err)
			}

			if c := string((contents)); c != tc.want {
				t.Fatalf("contents: got %v, want %v", c, tc.want)
			}
		})
	}
}
