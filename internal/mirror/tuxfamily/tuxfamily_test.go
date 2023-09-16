package tuxfamily

import (
	"errors"
	"reflect"
	"testing"

	"github.com/coffeebeats/gdenv/internal/mirror"
	"github.com/coffeebeats/gdenv/internal/version"
	"github.com/coffeebeats/gdenv/pkg/godot"
)

/* ----------------------- Test: TuxFamily.Executable ----------------------- */

func TestTuxFamilyExecutable(t *testing.T) {
	tests := []struct {
		ex   godot.Executable
		name string
		url  string
		err  error
	}{
		// Invalid inputs
		{ex: godot.Executable{}, err: mirror.ErrInvalidSpecification},
		{ex: godot.MustParseExecutable("Godot_v0.0.0-stable_linux.x86_64"), err: mirror.ErrInvalidSpecification},

		// Valid inputs
		{
			ex:   godot.MustParseExecutable("Godot_v4.1.1-stable_mono_linux_x86_64"),
			name: "Godot_v4.1.1-stable_mono_linux_x86_64.zip",
			url:  "https://downloads.tuxfamily.org/godotengine/4.1.1/mono/Godot_v4.1.1-stable_mono_linux_x86_64.zip",
		},
		{
			ex:   godot.MustParseExecutable("Godot_v4.1-stable_linux.x86_64"),
			name: "Godot_v4.1-stable_linux.x86_64.zip",
			url:  "https://downloads.tuxfamily.org/godotengine/4.1/Godot_v4.1-stable_linux.x86_64.zip",
		},
		{
			ex:   godot.MustParseExecutable("Godot_v4.0-dev.20220118_win64.exe"),
			name: "Godot_v4.0-dev.20220118_win64.exe.zip",
			url:  "https://downloads.tuxfamily.org/godotengine/4.0/pre-alpha/4.0-dev.20220118/Godot_v4.0-dev.20220118_win64.exe.zip",
		},
	}

	for _, tc := range tests {
		t.Run(tc.ex.String(), func(t *testing.T) {
			got, err := (&TuxFamily{}).Executable(tc.ex)

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

/* ------------------------ Test: TuxFamily.Checksum ------------------------ */

func TestTuxFamilyChecksum(t *testing.T) {
	tests := []struct {
		v    version.Version
		name string
		url  string
		err  error
	}{
		// Valid inputs
		{
			v:    version.MustParse("4.1.1-stable"),
			name: mirror.FilenameChecksums,
			url:  "https://downloads.tuxfamily.org/godotengine/4.1.1/" + mirror.FilenameChecksums,
		},
		{
			v:    version.MustParse("4.1-stable"),
			name: mirror.FilenameChecksums,
			url:  "https://downloads.tuxfamily.org/godotengine/4.1/" + mirror.FilenameChecksums,
		},
		{
			v:    version.MustParse("4.0-dev.20220118"),
			name: mirror.FilenameChecksums,
			url:  "https://downloads.tuxfamily.org/godotengine/4.0/pre-alpha/4.0-dev.20220118/" + mirror.FilenameChecksums,
		},
	}

	for _, tc := range tests {
		t.Run(tc.v.String(), func(t *testing.T) {
			got, err := (&TuxFamily{}).Checksum(tc.v)

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
