package checksum

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/coffeebeats/gdenv/internal/godot/artifact"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/source"
	"github.com/coffeebeats/gdenv/internal/godot/version"
)

/* ----------------------- Test: TestExtractExecutable ---------------------- */

func TestExtractExecutable(t *testing.T) {
	nameV4, nameV5 := "Godot_v4.0-stable_linux.x86_64", "Godot_v5.0-stable_linux.x86_64"

	archiveV4 := executable.Archive{Artifact: executable.MustParse(nameV4)}
	archiveV5 := executable.Archive{Artifact: executable.MustParse(nameV5)}

	tests := []struct {
		contents string
		exists   bool
		archive  executable.Archive
		want     string
		err      error
	}{
		// Invalid inputs
		{exists: false, archive: archiveV4, err: fs.ErrNotExist},
		{exists: true, contents: "abc 123 filename", archive: archiveV4, err: ErrUnrecognizedFormat},
		{
			exists:   true,
			contents: fmt.Sprintf("checksum %s", archiveV5.Name()),
			archive:  archiveV4,
			err:      ErrChecksumNotFound,
		},
		{
			exists:   true,
			contents: fmt.Sprintf("checksum1 %s\nchecksum2 %s", archiveV4.Name(), archiveV4.Name()),
			archive:  archiveV4,
			err:      ErrConflictingChecksum,
		},

		// Valid inputs
		{
			exists:   true,
			contents: fmt.Sprintf("checksumV4 %s\nchecksumV5 %s", archiveV4.Name(), archiveV5.Name()),
			archive:  archiveV5,
			want:     "checksumV5",
		},
	}

	for i, tc := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var c artifact.Local[Checksums[executable.Archive]]

			// NOTE: Property 'Artifact' doesn't need to be accessed.
			c.Path = filepath.Join(t.TempDir(), "checksums.txt")

			if tc.exists {
				if err := os.WriteFile(c.Path, []byte(tc.contents+"\n"), modeTestFile); err != nil {
					t.Fatalf("test setup: %#v", err)
				}
			}

			got, err := Extract(c, tc.archive)

			if !errors.Is(err, tc.err) {
				t.Errorf("err: got %#v, want %#v", err, tc.err)
			}

			if got != tc.want {
				t.Errorf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}

/* ------------------------- Test: TestExtractSource ------------------------ */

func TestExtractSource(t *testing.T) {
	sourceV3, sourceV4 := source.New(version.Godot3()), source.New(version.Godot4())

	archiveV3 := source.Archive{Artifact: sourceV3}
	archiveV4 := source.Archive{Artifact: sourceV4}

	tests := []struct {
		contents string
		exists   bool
		archive  source.Archive
		want     string
		err      error
	}{
		// Invalid inputs
		{exists: false, archive: archiveV4, err: fs.ErrNotExist},
		{exists: true, contents: "abc 123 filename", archive: archiveV4, err: ErrUnrecognizedFormat},
		{
			exists:   true,
			contents: fmt.Sprintf("checksum %s", archiveV4.Name()),
			archive:  archiveV3,
			err:      ErrChecksumNotFound,
		},
		{
			exists:   true,
			contents: fmt.Sprintf("checksum1 %s\nchecksum2 %s", archiveV3.Name(), archiveV3.Name()),
			archive:  archiveV3,
			err:      ErrConflictingChecksum,
		},

		// Valid inputs
		{
			exists:   true,
			contents: fmt.Sprintf("checksumV4 %s\nchecksumV5 %s", archiveV3.Name(), archiveV4.Name()),
			archive:  archiveV4,
			want:     "checksumV5",
		},
	}

	for i, tc := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var c artifact.Local[Checksums[source.Archive]]

			// NOTE: Property 'Artifact' doesn't need to be accessed.
			c.Path = filepath.Join(t.TempDir(), "checksums.txt")

			if tc.exists {
				if err := os.WriteFile(c.Path, []byte(tc.contents+"\n"), modeTestFile); err != nil {
					t.Fatalf("test setup: %#v", err)
				}
			}

			got, err := Extract(c, tc.archive)

			if !errors.Is(err, tc.err) {
				t.Errorf("err: got %#v, want %#v", err, tc.err)
			}

			if got != tc.want {
				t.Errorf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}
