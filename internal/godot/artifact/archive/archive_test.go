package archive

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

const (
	modeTestFile = 0600 // rw-------
)

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
				if err := os.WriteFile(archive.Path, []byte(""), modeTestFile); err != nil { // owner r+w
					t.Fatal(err)
				}

				// Given: The output path exists but is a file.
				if err := os.WriteFile(out, []byte(""), modeTestFile); err != nil { // owner r+w
					t.Fatal(err)
				}
			},

			err: fs.ErrInvalid,
		},

		// Valid inputs
		{
			name:    "archive fails to extract",
			archive: MockArchive[MockArtifact]{name: "archive.zip", err: fs.ErrPermission}, // arbitrary error
			out:     "directory",

			setUpFileSystem: func(t *testing.T, archive Local, out string) {
				// Given: The archive exists in the testing directory.
				if err := os.WriteFile(archive.Path, []byte(""), modeTestFile); err != nil { // owner r+w
					t.Fatal(err)
				}

				// Given: The output path does not exist.
			},

			err: fs.ErrPermission,
		},
		{
			name:    "archive extracts successfully with missing 'out'",
			archive: MockArchive[MockArtifact]{name: "archive.zip"},
			out:     "directory",

			setUpFileSystem: func(t *testing.T, archive Local, out string) {
				// Given: The archive exists in the testing directory.
				if err := os.WriteFile(archive.Path, []byte(""), modeTestFile); err != nil { // owner r+w
					t.Fatal(err)
				}

				// Given: The output path does not exist.
			},

			err: nil,
		},
		{
			name:    "archive extracts successfully with existing 'out' directory",
			archive: MockArchive[MockArtifact]{name: "archive.zip"},
			out:     "directory",

			setUpFileSystem: func(t *testing.T, archive Local, out string) {
				// Given: The archive exists in the testing directory.
				if err := os.WriteFile(archive.Path, []byte(""), modeTestFile); err != nil { // owner r+w
					t.Fatal(err)
				}

				// Given: The output path exists and is a directory.
				if err := os.MkdirAll(out, modeTestFile); err != nil { // owner r+w
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

			err := Extract(localArchive, dst)

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

var _ Archive = MockArchive[MockArtifact]{}

/* ----------------------------- Impl: Artifact ----------------------------- */

func (a MockArchive[T]) Name() string {
	return a.name
}

/* ------------------------------ Impl: Archive ----------------------------- */

func (a MockArchive[T]) extract(path, out string) error {
	return a.err
}

/* -------------------------------------------------------------------------- */
/*                            Struct: MockArtifact                            */
/* -------------------------------------------------------------------------- */

type MockArtifact struct {
	name string
}

/* ----------------------------- Impl: Artifact ----------------------------- */

func (a MockArtifact) Name() string {
	return a.name
}

/* ---------------------------- Impl: Archivable ---------------------------- */

func (a MockArtifact) Archivable() {}
