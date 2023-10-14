package store

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
		t.Fatalf("test setup: %v", err)
		return
	}

	tests := []struct {
		input string
		want  string
		err   error
	}{
		{"", "", ErrMissingPath},
		{"a", filepath.Join(wd, "a"), nil},
		{"a/b/c", filepath.Join(wd, "a/b/c"), nil},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got, err := Clean(tc.input)

			if !errors.Is(err, tc.err) {
				t.Errorf("err: got %#v, want %#v", err, tc.err)
			}
			if got != tc.want {
				t.Errorf("output: got %#v, want %#v", got, tc.want)
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
		{"a", true, true},
		{"a/b/c", true, true},
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
				store, err := Clean(path)
				if err != nil {
					t.Fatalf("test setup: %v", err)
				}

				if err := os.MkdirAll(store, modeTestDir); err != nil {
					t.Fatalf("test setup: %v", err)
				}
			}

			got := Exists(path)
			if got != tc.want {
				t.Errorf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}

/* ------------------------------- Test: Path ------------------------------- */

func TestPath(t *testing.T) {
	tests := []struct {
		env  string
		want string
		err  error
	}{
		{"", "", ErrMissingEnvVar},
		{"a", "", ErrInvalidPath},
		{"a/b/c", "", ErrInvalidPath},
		{"/", "/", nil},
		{"/a", "/a", nil},
		{"/a/b/c", "/a/b/c", nil},
	}

	for _, tc := range tests {
		t.Run(tc.env, func(t *testing.T) {
			t.Setenv(envVarStore, tc.env)

			got, err := Path()

			if !errors.Is(err, tc.err) {
				t.Errorf("err: got %#v, want %#v", err, tc.err)
			}

			if got != tc.want {
				t.Errorf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}
