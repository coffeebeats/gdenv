package github

import (
	"errors"
	"reflect"
	"testing"

	"github.com/coffeebeats/gdenv/internal/mirror"
	"github.com/coffeebeats/gdenv/pkg/godot"
)

/* ------------------------- Test: GitHub.Executable ------------------------ */

func TestGitHubExecutable(t *testing.T) {
	tests := []struct {
		ex   godot.Executable
		name string
		url  string
		err  error
	}{
		// Invalid inputs
		{ex: godot.Executable{}, err: mirror.ErrInvalidSpecification},
		{ex: godot.MustParseExecutable("Godot_v0.1.0-stable_linux.x86_64"), err: mirror.ErrInvalidSpecification},
		{ex: godot.MustParseExecutable("Godot_v4.1.1-unsupported-label_linux.x86_64"), err: mirror.ErrInvalidSpecification},

		// Valid inputs
		{
			ex:   godot.MustParseExecutable("Godot_v4.1.1-stable_linux.x86_64"),
			name: "Godot_v4.1.1-stable_linux.x86_64.zip",
			url:  "https://github.com/godotengine/godot/releases/download/4.1.1-stable/Godot_v4.1.1-stable_linux.x86_64.zip",
		},
	}

	for _, tc := range tests {
		t.Run(tc.ex.String(), func(t *testing.T) {
			got, err := (&GitHub{}).Executable(tc.ex)

			if !errors.Is(err, tc.err) {
				t.Fatalf("err: got %#v, want %#v", err, tc.err)
			}

			want, _ := mirror.NewAsset(tc.name, tc.url) // NOTE: Ignore 'err'; some expected.
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("output: got %#v, want %#v", got, want)
			}
		})
	}
}

/* -------------------------- Test: GitHub.Checksum ------------------------- */

func TestGitHubChecksum(t *testing.T) {
	tests := []struct {
		v    godot.Version
		name string
		url  string
		err  error
	}{
		// Invalid inputs
		{v: godot.Version{}, err: mirror.ErrInvalidSpecification},
		{v: godot.MustParseVersion("v0.0.0"), err: mirror.ErrInvalidSpecification},
		{v: godot.MustParseVersion("v4.1.1-unsupported-label"), err: mirror.ErrInvalidSpecification},

		// Valid inputs
		{
			v:    godot.MustParseVersion("4.1.1-stable"),
			name: mirror.FilenameChecksums,
			url:  "https://github.com/godotengine/godot/releases/download/4.1.1-stable/" + mirror.FilenameChecksums,
		},
	}

	for _, tc := range tests {
		t.Run(tc.v.String(), func(t *testing.T) {
			got, err := (&GitHub{}).Checksum(tc.v)

			if !errors.Is(err, tc.err) {
				t.Fatalf("err: got %#v, want %#v", err, tc.err)
			}

			want, _ := mirror.NewAsset(tc.name, tc.url) // NOTE: Ignore 'err'; some expected.
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("output: got %#v, want %#v", got, want)
			}
		})
	}
}
