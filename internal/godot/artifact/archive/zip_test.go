package archive

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/coffeebeats/gdenv/internal/fstest"
	"github.com/coffeebeats/gdenv/internal/osutil"
	"github.com/coffeebeats/gdenv/internal/progress"
)

/* -------------------------- Function: TestExtract ------------------------- */

func TestZipExtract(t *testing.T) {
	tests := []struct {
		name     string
		artifact fstest.Writer

		want []fstest.Asserter
		err  error
	}{
		// Invalid inputs
		{
			name:     "missing archive returns file not found",
			artifact: fstest.Absent{Path: "archive.zip"},

			err: os.ErrNotExist,
		},

		// Valid inputs
		{
			name:     "empty archive doesn't return an error",
			artifact: fstest.Zip{Path: "archive.zip"},
		},
		{
			name: "multiple files can be extracted",
			artifact: fstest.Zip{
				Path: "archive.zip",
				// Relative to archive file.
				Contents: []fstest.Writer{
					fstest.File{Path: "godot1.exe", Contents: "1"},
					fstest.File{Path: "godot2.exe", Contents: "2"},
				},
			},

			// Relative to extraction directory.
			want: []fstest.Asserter{
				fstest.File{Path: "godot1.exe", Contents: "1"},
				fstest.File{Path: "godot2.exe", Contents: "2"},
			},
		},
		{
			name: "nested files can be extracted",
			artifact: fstest.Zip{
				Path: "archive.zip",
				// Relative to archive file.
				Contents: []fstest.Writer{
					fstest.Dir{Path: "godot"},
					fstest.File{Path: "godot/godot1.exe", Contents: "1"},
				},
			},

			// Relative to extraction directory.
			want: []fstest.Asserter{
				fstest.Dir{Path: "godot"},
				fstest.File{Path: "godot/godot1.exe", Contents: "1"},
			},
		},
		{
			name: "nested files without parent directory can be extracted",
			artifact: fstest.Zip{
				Path: "archive.zip",
				// Relative to archive file.
				Contents: []fstest.Writer{
					fstest.File{Path: "godot/godot1.exe", Contents: "1"},
				},
			},

			// Relative to extraction directory.
			want: []fstest.Asserter{
				fstest.Dir{Path: "godot"},
				fstest.File{Path: "godot/godot1.exe", Contents: "1"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmp := t.TempDir()

			// Given: A directory to extract into.
			out := filepath.Join(tmp, "extract")
			if err := os.Mkdir(out, osutil.ModeUserRWX); err != nil {
				t.Fatal(err)
			}

			// Given: The specified archive exists on the file system.
			tc.artifact.Write(t, tmp)

			// When: The archive is extracted.
			err := (Zip[MockArtifact]{}).extract(context.Background(), tc.artifact.Abs(t, tmp), out)

			// Then: The expected error value is returned.
			if !errors.Is(err, tc.err) {
				t.Errorf("got: %v, want: %v", err, tc.err)
			}

			// Then: The expected files exist on the file system.
			for _, f := range tc.want {
				f.Assert(t, out)
			}
		})
	}
}

/* -------------------- Function: TestExtractWithProgress ------------------- */

func TestZipExtractWithProgress(t *testing.T) {
	tests := []struct {
		name     string
		artifact fstest.Writer
		progress *progress.Progress

		want float64
		err  error
	}{}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmp := t.TempDir()

			// Given: A directory to extract into.
			out := filepath.Join(tmp, "extract")

			// Given: The specified archive exists on the file system.
			tc.artifact.Write(t, tmp)

			// Given: A 'context.Context' with the specified progress reporter.
			ctx := context.Background()
			if tc.progress != nil {
				ctx = WithProgress(ctx, tc.progress)
			}

			// When: The archive is extracted.
			// Then: The expected error value is returned.
			if err := (Zip[MockArtifact]{}).extract(ctx, tc.artifact.Abs(t, tmp), out); !errors.Is(err, tc.err) {
				t.Errorf("got: %v, want: %v", err, tc.err)
			}

			if tc.progress == nil {
				return
			}

			// Then: The progress value should be 100%.
			if got := tc.progress.Percentage(); got != tc.want {
				t.Errorf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}
