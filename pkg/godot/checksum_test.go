package godot

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

/* ------------------------ Test: TestComputeChecksum ----------------------- */

func TestComputeChecksum(t *testing.T) {
	tests := []struct {
		f        string
		contents string
		exists   bool
		want     string
		err      error
	}{
		// Invalid inputs
		{f: "", exists: false, err: ErrMissingPath},
		{f: "executable.zip", exists: false, err: ErrFileSystem},

		// Valid inputs
		{
			f:        "executable.zip",
			contents: "abc",
			exists:   true,
			want:     "4f285d0c0cc77286d8731798b7aae2639e28270d4166f40d769cbbdca5230714d848483d364e2f39fe6cb9083c15229b39a33615ebc6d57605f7c43f6906739d",
			err:      nil,
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d-'%s'-%v", i, tc.f, tc.exists), func(t *testing.T) {
			if tc.exists {
				tmp := t.TempDir()
				tc.f = filepath.Join(tmp, tc.f)

				if err := os.WriteFile(tc.f, []byte(tc.contents+"\n"), os.ModePerm); err != nil {
					t.Fatalf("test setup: %#v", err)
				}
			}

			got, err := ComputeChecksum(tc.f)

			if !errors.Is(err, tc.err) {
				t.Fatalf("err: got %#v, want %#v", err, tc.err)
			}

			if got != tc.want {
				t.Fatalf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}

/* ------------------------ Test: TestExtractChecksum ----------------------- */

func TestExtractChecksum(t *testing.T) {
	nameV4, nameV5 := "Godot_v4.0-stable_linux.x86_64", "Godot_v5.0-stable_linux.x86_64"
	exV4, exV5 := MustParseExecutable(nameV4), MustParseExecutable(nameV5)

	tests := []struct {
		f        string
		contents string
		exists   bool
		ex       Executable
		want     string
		err      error
	}{
		// Invalid inputs
		{f: "", exists: false, err: ErrMissingPath},
		{f: "checksums.txt", exists: false, ex: exV4, err: ErrFileSystem},
		{f: "checksums.txt", exists: true, contents: "abc 123 filename", ex: exV4, err: ErrUnrecognizedFormat},
		{f: "checksums.txt", exists: true, contents: fmt.Sprintf("checksum %s", nameV5), ex: exV4, err: ErrMissingChecksum},
		{f: "checksums.txt", exists: true, contents: fmt.Sprintf("checksum1 %s\nchecksum2 %s", nameV4, nameV4), ex: exV4, err: ErrConflictingChecksum},

		// Valid inputs
		{
			f:        "checksums.txt",
			contents: fmt.Sprintf("checksumV4 %s\nchecksumV5 %s", nameV4, nameV5),
			exists:   true,
			ex:       exV5,
			want:     "checksumV5",
			err:      nil,
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d-'%s'-%v", i, tc.f, tc.exists), func(t *testing.T) {
			if tc.exists {
				tmp := t.TempDir()
				tc.f = filepath.Join(tmp, tc.f)

				if err := os.WriteFile(tc.f, []byte(tc.contents+"\n"), os.ModePerm); err != nil {
					t.Fatalf("test setup: %#v", err)
				}
			}

			got, err := ExtractChecksum(tc.f, tc.ex)

			if !errors.Is(err, tc.err) {
				t.Fatalf("err: got %#v, want %#v", err, tc.err)
			}

			if got != tc.want {
				t.Fatalf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}
