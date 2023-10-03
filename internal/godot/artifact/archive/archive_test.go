package archive

import (
	"errors"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"

	"github.com/coffeebeats/gdenv/internal/godot/artifact"
)

/* ------------------------------ Test: Extract ----------------------------- */

func TestExtract(t *testing.T) {
	tests := []struct {
		archive Zip[testArtifact]
		out     string

		want artifact.Local[testArtifact]
		err  error
	}{
		// Valid inputs
		{},

		// Invalid inputs
		{},
	}

	for i, tc := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			tmp := t.TempDir()

			err := Extract(
				Local{
					Artifact: tc.archive,
					Path:     filepath.Join(tmp, tc.archive.Name()),
				},
				tc.out,
			)

			if !errors.Is(err, tc.err) {
				t.Errorf("err: got %v, want %v", err, tc.err)
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("output: got %v, want %v", got, tc.want)
			}
		})
	}
}

/* -------------------------------------------------------------------------- */
/*                            Struct: testArtifact                            */
/* -------------------------------------------------------------------------- */

type testArtifact struct{}

/* ----------------------------- Impl: Artifact ----------------------------- */

func (t testArtifact) Name() string {
	return ""
}

/* ---------------------------- Impl: Archivable ---------------------------- */

func (t testArtifact) Archivable() {}
