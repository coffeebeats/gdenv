package checksum

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/coffeebeats/gdenv/internal/godot/artifact"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/executable"
)

/* ---------------------------- Test: TestCompute --------------------------- */

func TestCompute(t *testing.T) {
	tests := []struct {
		contents string
		exists   bool
		want     string
		err      error
	}{
		// Invalid inputs
		{exists: false, err: fs.ErrNotExist},

		// Valid inputs
		{
			contents: "abc",
			exists:   true,
			want:     "4f285d0c0cc77286d8731798b7aae2639e28270d4166f40d769cbbdca5230714d848483d364e2f39fe6cb9083c15229b39a33615ebc6d57605f7c43f6906739d",
			err:      nil,
		},
		// TODO: Add a test case for a 'source.Archive'.
	}

	for i, tc := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var f artifact.Local[executable.Archive]

			// NOTE: Property 'Artifact' doesn't need to be accessed.
			f.Path = filepath.Join(t.TempDir(), "archive.zip")

			if tc.exists {
				if err := os.WriteFile(f.Path, []byte(tc.contents+"\n"), os.ModePerm); err != nil {
					t.Fatalf("test setup: %#v", err)
				}
			}

			got, err := Compute(f)

			if !errors.Is(err, tc.err) {
				t.Errorf("err: got %#v, want %#v", err, tc.err)
			}

			if got != tc.want {
				t.Errorf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}
