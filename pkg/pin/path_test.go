package pin

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

/* ------------------------------- Test: Clean ------------------------------ */

func TestClean(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("%v", err)
		return
	}

	tests := []struct {
		input string
		want  string
		err   error
	}{
		{"", "", ErrMissingPath},
		{"a/b/c", filepath.Join(wd, "a/b/c", pinFilename), nil},
		{"a/b/c/" + pinFilename, filepath.Join(wd, "a/b/c", pinFilename), nil},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got, err := Clean(tc.input)

			if !errors.Is(err, tc.err) {
				t.Fatalf("err: got %#v, want %#v", err, tc.err)
			}
			if got != tc.want {
				t.Fatalf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}

/* ------------------------------ Test: Exists ------------------------------ */

func TestExists(t *testing.T) {
	tests := []struct {
		path       string
		isRelative bool // Is the path relative to 'tmp'
		want       bool
	}{
		{"", true, true},
		{"a", true, true},
		{"a/b/c", true, true},
		{"", true, false},
		{"a", true, false},
		{"a/b/c", true, false},

		// Check the empty string
		{"", false, false},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			tmp := t.TempDir()

			path := filepath.Join(tmp, tc.path)
			if !tc.isRelative {
				path = tc.path
			}

			// Create the pin file
			if tc.want {
				pin, err := Clean(path)
				if err != nil {
					t.Fatalf("test setup: %v", err)
				}

				if err := os.MkdirAll(filepath.Dir(pin), os.ModePerm); err != nil {
					t.Fatalf("test setup: %v", err)
				}

				if err := os.WriteFile(pin, []byte(""), os.ModePerm); err != nil {
					t.Fatalf("test setup: %v", err)
				}
			}

			got := Exists(path)
			if got != tc.want {
				t.Fatalf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}
