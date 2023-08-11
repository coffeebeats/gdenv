package pin

import (
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
			got, _ := Clean(tc.input)

			// if err != tc.err && errors.Unwrap(err) != tc.err {
			// 	t.Fatalf("err: got %#v, want %#v", err, tc.err)
			// }
			if got != tc.want {
				t.Fatalf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}
