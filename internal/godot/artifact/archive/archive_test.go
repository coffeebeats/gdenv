package archive

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/coffeebeats/gdenv/internal/godot/artifact/artifacttest"
	"github.com/coffeebeats/gdenv/internal/osutil"
)

type MockArtifact = artifacttest.MockArtifact

/* ------------------------------ Test: Extract ----------------------------- */

func TestExtract(t *testing.T) {
	tests := []struct {
		name    string
		archive MockArchive[MockArtifact]
		out     string

		setUpFileSystem func(*testing.T, Local, string)

		err error
	}{
		// Invalid inputs
		{
			name:    "input artifact does not exist",
			archive: MockArchive[MockArtifact]{name: "archive.zip"},
			out:     "name",

			setUpFileSystem: func(t *testing.T, archive Local, out string) {
				// Given: The archive does not exist.
				// Given: The output path does not exist.
			},

			err: fs.ErrNotExist,
		},
		{
			name:    "'out' path exists but is a file",
			archive: MockArchive[MockArtifact]{name: "archive.zip"},
			out:     "directory",

			setUpFileSystem: func(t *testing.T, archive Local, out string) {
				// Given: The archive exists in the testing directory.
				if err := os.WriteFile(archive.Path, []byte(""), osutil.ModeUserRW); err != nil { // owner r+w
					t.Fatal(err)
				}

				// Given: The output path exists but is a file.
				if err := os.WriteFile(out, []byte(""), osutil.ModeUserRW); err != nil { // owner r+w
					t.Fatal(err)
				}
			},

			err: fs.ErrInvalid,
		},
		{
			name:    "fails with missing 'out'",
			archive: MockArchive[MockArtifact]{name: "archive.zip"},
			out:     "directory",

			setUpFileSystem: func(t *testing.T, archive Local, out string) {
				// Given: The archive exists in the testing directory.
				if err := os.WriteFile(archive.Path, []byte(""), osutil.ModeUserRW); err != nil { // owner r+w
					t.Fatal(err)
				}

				// Given: The output path does not exist.
			},

			err: fs.ErrNotExist,
		},
		// Valid inputs
		{
			name:    "archive extracts successfully with existing 'out' directory",
			archive: MockArchive[MockArtifact]{name: "archive.zip"},
			out:     "directory",

			setUpFileSystem: func(t *testing.T, archive Local, out string) {
				// Given: The archive exists in the testing directory.
				if err := os.WriteFile(archive.Path, []byte(""), osutil.ModeUserRW); err != nil { // owner r+w
					t.Fatal(err)
				}

				// Given: The output path exists and is a directory.
				if err := os.MkdirAll(out, osutil.ModeUserRWX); err != nil { // owner r+w
					t.Fatal(err)
				}
			},

			err: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmp := t.TempDir()

			src := filepath.Join(tmp, tc.archive.Name())
			dst := filepath.Join(tmp, tc.out)

			localArchive := Local{
				Artifact: tc.archive,
				Path:     src,
			}

			tc.setUpFileSystem(t, localArchive, dst)

			err := Extract(context.Background(), localArchive, dst)

			if !errors.Is(err, tc.err) {
				t.Errorf("err: got %v, want %v", err, tc.err)
			}
		})
	}
}

/* -------------------------------------------------------------------------- */
/*                             Struct: MockArchive                            */
/* -------------------------------------------------------------------------- */

type MockArchive[T Archivable] struct {
	name string
	err  error
}

var _ Archive = MockArchive[artifacttest.MockArtifact]{}

/* ----------------------------- Impl: Artifact ----------------------------- */

func (a MockArchive[T]) Name() string {
	return a.name
}

/* ------------------------------ Impl: Archive ----------------------------- */

func (a MockArchive[T]) extract(_ context.Context, path, out string) error {
	return a.err
}
