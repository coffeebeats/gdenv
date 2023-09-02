package mirror

import (
	"errors"
	"net/url"
	"reflect"
	"testing"

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
		{ex: godot.Executable{}, err: ErrInvalidSpecification},
		{ex: godot.MustParseExecutable("Godot_v0.1.0-stable_linux.x86_64"), err: ErrInvalidSpecification},
		{ex: godot.MustParseExecutable("Godot_v4.1.1-unsupported-label_linux.x86_64"), err: ErrInvalidSpecification},

		// Valid inputs
		{
			ex:   godot.MustParseExecutable("Godot_v4.1.1-stable_linux.x86_64"),
			name: "Godot_v4.1.1-stable_linux.x86_64.zip",
			url:  "https://github.com/godotengine/godot/releases/download/4.1.1-stable/Godot_v4.1.1-stable_linux.x86_64.zip",
		},
	}

	for _, tc := range tests {
		t.Run(tc.ex.String(), func(t *testing.T) {
			u, err := url.Parse(tc.url)
			if err != nil {
				t.Fatalf("test setup: %#v", err)
			}
			if tc.url == "" {
				u = nil
			}

			got, err := (&GitHub{}).Executable(tc.ex)

			if !errors.Is(err, tc.err) {
				t.Fatalf("err: got %#v, want %#v", err, tc.err)
			}

			want := Asset{name: tc.name, url: u}
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
		{v: godot.Version{}, err: ErrInvalidSpecification},
		{v: godot.MustParseVersion("v0.0.0"), err: ErrInvalidSpecification},
		{v: godot.MustParseVersion("v4.1.1-unsupported-label"), err: ErrInvalidSpecification},

		// Valid inputs
		{
			v:    godot.MustParseVersion("4.1.1-stable"),
			name: filenameChecksums,
			url:  "https://github.com/godotengine/godot/releases/download/4.1.1-stable/" + filenameChecksums,
		},
	}

	for _, tc := range tests {
		t.Run(tc.v.String(), func(t *testing.T) {
			u, err := url.Parse(tc.url)
			if err != nil {
				t.Fatalf("test setup: %#v", err)
			}
			if tc.url == "" {
				u = nil
			}

			got, err := (&GitHub{}).Checksum(tc.v)

			if !errors.Is(err, tc.err) {
				t.Fatalf("err: got %#v, want %#v", err, tc.err)
			}

			want := Asset{name: tc.name, url: u}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("output: got %#v, want %#v", got, want)
			}
		})
	}
}
