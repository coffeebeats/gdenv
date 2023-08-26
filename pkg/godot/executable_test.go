package godot

import (
	"errors"
	"fmt"
	"testing"
)

/* -------------------------- Test: Executable.Name ------------------------- */

func TestExecutableName(t *testing.T) {
	tests := []struct {
		p    Platform
		v    Version
		want string
		err  error
	}{
		// Invalid inputs
		{Platform{}, Version{}, "", ErrMissingOS},
		{Platform{OS: linux}, Version{}, "", ErrMissingArch},
		{Platform{OS: linux, Arch: amd64}, Version{}, "", ErrUnsupportedVersion},

		{Platform{OS: linux, Arch: amd64}, Version{major: 2}, "", ErrUnsupportedVersion},
		{Platform{OS: macOS, Arch: amd64}, Version{major: 2}, "", ErrUnsupportedVersion},
		{Platform{OS: windows, Arch: amd64}, Version{major: 2}, "", ErrUnsupportedVersion},

		// Valid inputs - Linux

		// v3.0
		{Platform{OS: linux, Arch: i386}, Version{major: 3}, "Godot_v3.0-stable_x11.32", nil},
		{Platform{OS: linux, Arch: amd64}, Version{major: 3}, "Godot_v3.0-stable_x11.64", nil},
		{Platform{OS: linux, Arch: arm64}, Version{major: 3}, "", ErrUnsupportedArch},
		{Platform{OS: linux, Arch: universal}, Version{major: 3}, "", ErrUnsupportedArch},

		// v3.6-beta2
		{Platform{OS: linux, Arch: i386}, Version{major: 3, minor: 6, label: "beta2"}, "Godot_v3.6-beta2_x11.32", nil},
		{Platform{OS: linux, Arch: amd64}, Version{major: 3, minor: 6, label: "beta2"}, "Godot_v3.6-beta2_x11.64", nil},
		{Platform{OS: linux, Arch: arm64}, Version{major: 3, minor: 6, label: "beta2"}, "", ErrUnsupportedArch},
		{Platform{OS: linux, Arch: universal}, Version{major: 3, minor: 6, label: "beta2"}, "", ErrUnsupportedArch},

		// v4.0-rc4
		{Platform{OS: linux, Arch: i386}, Version{major: 4, label: "rc4"}, "Godot_v4.0-rc4_linux.x86_32", nil},
		{Platform{OS: linux, Arch: amd64}, Version{major: 4, label: "rc4"}, "Godot_v4.0-rc4_linux.x86_64", nil},
		{Platform{OS: linux, Arch: arm64}, Version{major: 4, label: "rc4"}, "", ErrUnsupportedArch},
		{Platform{OS: linux, Arch: universal}, Version{major: 4, label: "rc4"}, "", ErrUnsupportedArch},

		// Valid inputs - MacOS

		// v3.0
		{Platform{OS: macOS, Arch: i386}, Version{major: 3}, "Godot_v3.0-stable_osx.fat", nil},
		{Platform{OS: macOS, Arch: amd64}, Version{major: 3}, "Godot_v3.0-stable_osx.fat", nil},
		{Platform{OS: macOS, Arch: arm64}, Version{major: 3}, "", ErrUnsupportedArch},
		{Platform{OS: macOS, Arch: universal}, Version{major: 3}, "", ErrUnsupportedArch},

		// v3.6-beta2
		{Platform{OS: macOS, Arch: amd64}, Version{major: 3, minor: 6, label: "beta2"}, "Godot_v3.6-beta2_osx.universal", nil},
		{Platform{OS: macOS, Arch: arm64}, Version{major: 3, minor: 6, label: "beta2"}, "Godot_v3.6-beta2_osx.universal", nil},
		{Platform{OS: macOS, Arch: i386}, Version{major: 3, minor: 6, label: "beta2"}, "", ErrUnsupportedArch},
		{Platform{OS: macOS, Arch: universal}, Version{major: 3, minor: 6, label: "beta2"}, "", ErrUnsupportedArch},

		// v4.0-rc4
		{Platform{OS: macOS, Arch: amd64}, Version{major: 4, label: "rc4"}, "Godot_v4.0-rc4_macos.universal", nil},
		{Platform{OS: macOS, Arch: arm64}, Version{major: 4, label: "rc4"}, "Godot_v4.0-rc4_macos.universal", nil},
		{Platform{OS: macOS, Arch: i386}, Version{major: 4, label: "rc4"}, "", ErrUnsupportedArch},
		{Platform{OS: macOS, Arch: universal}, Version{major: 4, label: "rc4"}, "", ErrUnsupportedArch},

		// Valid inputs - Windows

		// v3.0
		{Platform{OS: windows, Arch: i386}, Version{major: 3}, "Godot_v3.0-stable_win32.exe", nil},
		{Platform{OS: windows, Arch: amd64}, Version{major: 3}, "Godot_v3.0-stable_win64.exe", nil},
		{Platform{OS: windows, Arch: arm64}, Version{major: 3}, "", ErrUnsupportedArch},
		{Platform{OS: windows, Arch: universal}, Version{major: 3}, "", ErrUnsupportedArch},

		// v3.6-beta2
		{Platform{OS: windows, Arch: i386}, Version{major: 3, minor: 6, label: "beta2"}, "Godot_v3.6-beta2_win32.exe", nil},
		{Platform{OS: windows, Arch: amd64}, Version{major: 3, minor: 6, label: "beta2"}, "Godot_v3.6-beta2_win64.exe", nil},
		{Platform{OS: windows, Arch: arm64}, Version{major: 3, minor: 6, label: "beta2"}, "", ErrUnsupportedArch},
		{Platform{OS: windows, Arch: universal}, Version{major: 3, minor: 6, label: "beta2"}, "", ErrUnsupportedArch},

		// v4.0-rc4
		{Platform{OS: windows, Arch: i386}, Version{major: 4, label: "rc4"}, "Godot_v4.0-rc4_win32.exe", nil},
		{Platform{OS: windows, Arch: amd64}, Version{major: 4, label: "rc4"}, "Godot_v4.0-rc4_win64.exe", nil},
		{Platform{OS: windows, Arch: arm64}, Version{major: 4, label: "rc4"}, "", ErrUnsupportedArch},
		{Platform{OS: windows, Arch: universal}, Version{major: 4, label: "rc4"}, "", ErrUnsupportedArch},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got, err := Executable{tc.p, tc.v}.Name()

			if !errors.Is(err, tc.err) {
				t.Fatalf("err: got %#v, want %#v", err, tc.err)
			}
			if got != tc.want {
				t.Fatalf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}

/* -------------------------- Test: ParseExecutable ------------------------- */

func TestParseExecutable(t *testing.T) {
	tests := []struct {
		s    string
		want Executable
		err  error
	}{
		// Invalid inputs
		{s: "", want: Executable{}, err: ErrMissingName},

		// Valid inputs

	}

	for i, tc := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got, err := ParseExecutable(tc.s)

			if !errors.Is(err, tc.err) {
				t.Fatalf("err: got %#v, want %#v", err, tc.err)
			}
			if got != tc.want {
				t.Fatalf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}
