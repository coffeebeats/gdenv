package pin

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

/* ------------------------------- Test: clean ------------------------------ */

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
		{"a/b/c", filepath.Join(wd, "a/b/c", pinFilename), nil},
		{"a/b/c/" + pinFilename, filepath.Join(wd, "a/b/c", pinFilename), nil},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got, err := clean(tc.input)

			if !errors.Is(err, tc.err) {
				t.Errorf("err: got %#v, want %#v", err, tc.err)
			}
			if got != tc.want {
				t.Errorf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}
