package mirror

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"testing"

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
		{ex: godot.Executable{}, err: ErrInvalidSpecification},
		{ex: godot.MustParseExecutable("Godot_v0.0.0-stable_linux.x86_64"), err: ErrInvalidSpecification},

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
			u, err := url.Parse(tc.url)
			if err != nil {
				t.Fatalf("test setup: %#v", err)
			}
			if tc.url == "" {
				u = nil
			}

			got, err := (&TuxFamily{}).Executable(tc.ex)

			if !errors.Is(err, tc.err) {
				t.Fatalf("err: got %#v, want %#v", err, tc.err)
			}

			want := Asset{name: tc.name, url: u}
			if !reflect.DeepEqual(got, want) {
				fmt.Println(got.URL().String())
				t.Fatalf("output: got %#v, want %#v", got, want)
			}
		})
	}
}

/* ------------------------ Test: TuxFamily.Checksum ------------------------ */

func TestTuxFamilyChecksum(t *testing.T) {
	tests := []struct {
		v    godot.Version
		name string
		url  string
		err  error
	}{
		// Valid inputs
		{
			v:    godot.MustParseVersion("4.1.1-stable"),
			name: filenameChecksums,
			url:  "https://downloads.tuxfamily.org/godotengine/4.1.1/" + filenameChecksums,
		},
		{
			v:    godot.MustParseVersion("4.1-stable"),
			name: filenameChecksums,
			url:  "https://downloads.tuxfamily.org/godotengine/4.1/" + filenameChecksums,
		},
		{
			v:    godot.MustParseVersion("4.0-dev.20220118"),
			name: filenameChecksums,
			url:  "https://downloads.tuxfamily.org/godotengine/4.0/pre-alpha/4.0-dev.20220118/" + filenameChecksums,
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

			got, err := (&TuxFamily{}).Checksum(tc.v)

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
