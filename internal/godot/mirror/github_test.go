package mirror

import (
	"errors"
	"net/url"
	"reflect"
	"testing"

	"github.com/coffeebeats/gdenv/internal/godot/artifact"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/checksum"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/internal/godot/version"
)

/* ------------------------- Test: GitHub.Executable ------------------------ */

func TestGitHubExecutable(t *testing.T) {
	tests := []struct {
		ex   executable.Executable
		name string
		url  *url.URL
		err  error
	}{
		// Invalid inputs
		{ex: executable.Executable{}, err: ErrInvalidSpecification},
		{ex: executable.MustParse("Godot_v0.1.0-stable_linux.x86_64"), err: ErrInvalidSpecification},
		{ex: executable.MustParse("Godot_v4.1.1-unsupported-label_linux.x86_64"), err: ErrInvalidSpecification},

		// Valid inputs
		{
			ex:   executable.MustParse("Godot_v4.1.1-stable_linux.x86_64"),
			name: "Godot_v4.1.1-stable_linux.x86_64.zip",
			url:  mustParseURL(t, "https://github.com/godotengine/godot/releases/download/4.1.1-stable/Godot_v4.1.1-stable_linux.x86_64.zip"),
		},
		{
			ex:   executable.MustParse("Godot_v4.1-stable_linux.x86_64"),
			name: "Godot_v4.1-stable_linux.x86_64.zip",
			url:  mustParseURL(t, "https://github.com/godotengine/godot/releases/download/4.1-stable/Godot_v4.1-stable_linux.x86_64.zip"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.ex.String(), func(t *testing.T) {
			got, err := (&GitHub{}).ExecutableArchive(tc.ex)

			if !errors.Is(err, tc.err) {
				t.Errorf("err: got %v, want %v", err, tc.err)
			}

			if got := got.Artifact.Name(); got != tc.name {
				t.Errorf("output: got %v, want %v", got, tc.name)
			}
			if got := got.URL; !reflect.DeepEqual(got, tc.url) {
				t.Errorf("output: got %v, want %v", got, tc.url)
			}
		})
	}
}

/* -------------------------- Test: GitHub.Checksum ------------------------- */

func TestGitHubChecksum(t *testing.T) {
	tests := []struct {
		v   version.Version
		url *url.URL
		err error
	}{
		// Invalid inputs
		{v: version.Version{}, err: ErrInvalidSpecification},
		{v: version.MustParse("v0.0.0"), err: ErrInvalidSpecification},
		{v: version.MustParse("v4.1.1-unsupported-label"), err: ErrInvalidSpecification},

		// Valid inputs
		{
			v:   version.MustParse("4.1.1-stable"),
			url: mustParseURL(t, "https://github.com/godotengine/godot/releases/download/4.1.1-stable/SHA512-SUMS.txt"),
		},
		{
			v:   version.MustParse("4.1.0-stable"),
			url: mustParseURL(t, "https://github.com/godotengine/godot/releases/download/4.1-stable/SHA512-SUMS.txt"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.v.String(), func(t *testing.T) {
			got, err := (&GitHub{}).ExecutableArchiveChecksums(tc.v)

			if !errors.Is(err, tc.err) {
				t.Errorf("err: got %#v, want %#v", err, tc.err)
			}

			// The test setup below will fail for invalid inputs.
			if tc.url == nil {
				return
			}

			ex, err := checksum.NewExecutable(tc.v)
			if err != nil {
				t.Fatalf("test setup: %v", err)
			}

			want := artifact.Remote[checksum.Executable]{Artifact: ex, URL: tc.url}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("output: got %#v, want %#v", got, want)
			}
		})
	}
}
