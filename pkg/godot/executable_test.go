package godot

import (
	"errors"
	"fmt"
	"testing"
)

/* -------------------------- Test: Executable.Name ------------------------- */

func TestExecutableName(t *testing.T) {
	var (
		v3     = Version{3, 6, 0, "beta1"}
		v4     = Version{major: 4}
		v4Mono = Version{major: 4, label: "stable_mono"}
		v5     = Version{5, 0, 0, "rc4"}
	)

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

		// v3.6-beta1
		{Platform{OS: linux, Arch: i386}, v3, "Godot_v3.6-beta1_x11.32", nil},
		{Platform{OS: linux, Arch: amd64}, v3, "Godot_v3.6-beta1_x11.64", nil},
		{Platform{OS: linux, Arch: arm64}, v3, "", ErrUnsupportedArch},
		{Platform{OS: linux, Arch: universal}, v3, "", ErrUnsupportedArch},

		// v4.0
		{Platform{OS: linux, Arch: i386}, v4, "Godot_v4.0-stable_linux.x86_32", nil},
		{Platform{OS: linux, Arch: amd64}, v4, "Godot_v4.0-stable_linux.x86_64", nil},
		{Platform{OS: linux, Arch: arm64}, v4, "", ErrUnsupportedArch},
		{Platform{OS: linux, Arch: universal}, v4, "", ErrUnsupportedArch},

		// v4.0-stable_mono
		{Platform{OS: linux, Arch: i386}, v4Mono, "Godot_v4.0-stable_mono_linux_x86_32", nil},
		{Platform{OS: linux, Arch: amd64}, v4Mono, "Godot_v4.0-stable_mono_linux_x86_64", nil},
		{Platform{OS: linux, Arch: arm64}, v4Mono, "", ErrUnsupportedArch},
		{Platform{OS: linux, Arch: universal}, v4Mono, "", ErrUnsupportedArch},

		// v5.0-rc4
		{Platform{OS: linux, Arch: i386}, v5, "Godot_v5.0-rc4_linux.x86_32", nil},
		{Platform{OS: linux, Arch: amd64}, v5, "Godot_v5.0-rc4_linux.x86_64", nil},
		{Platform{OS: linux, Arch: arm64}, v5, "", ErrUnsupportedArch},
		{Platform{OS: linux, Arch: universal}, v5, "", ErrUnsupportedArch},

		// Valid inputs - MacOS

		// v3.6-beta1
		{Platform{OS: macOS, Arch: amd64}, v3, "Godot_v3.6-beta1_osx.universal", nil},
		{Platform{OS: macOS, Arch: arm64}, v3, "Godot_v3.6-beta1_osx.universal", nil},
		{Platform{OS: macOS, Arch: i386}, v3, "", ErrUnsupportedArch},
		{Platform{OS: macOS, Arch: universal}, v3, "", ErrUnsupportedArch},

		// v4.0
		{Platform{OS: macOS, Arch: amd64}, v4, "Godot_v4.0-stable_macos.universal", nil},
		{Platform{OS: macOS, Arch: arm64}, v4, "Godot_v4.0-stable_macos.universal", nil},
		{Platform{OS: macOS, Arch: i386}, v4, "", ErrUnsupportedArch},
		{Platform{OS: macOS, Arch: universal}, v4, "", ErrUnsupportedArch},

		// v4.0-stable_mono
		{Platform{OS: macOS, Arch: amd64}, v4Mono, "Godot_v4.0-stable_mono_macos.universal", nil},
		{Platform{OS: macOS, Arch: arm64}, v4Mono, "Godot_v4.0-stable_mono_macos.universal", nil},
		{Platform{OS: macOS, Arch: i386}, v4Mono, "", ErrUnsupportedArch},
		{Platform{OS: macOS, Arch: universal}, v4Mono, "", ErrUnsupportedArch},

		// v5.0-rc4
		{Platform{OS: macOS, Arch: amd64}, v5, "Godot_v5.0-rc4_macos.universal", nil},
		{Platform{OS: macOS, Arch: arm64}, v5, "Godot_v5.0-rc4_macos.universal", nil},
		{Platform{OS: macOS, Arch: i386}, v5, "", ErrUnsupportedArch},
		{Platform{OS: macOS, Arch: universal}, v5, "", ErrUnsupportedArch},

		// Valid inputs - Windows

		// v3.6-beta1
		{Platform{OS: windows, Arch: i386}, v3, "Godot_v3.6-beta1_win32.exe", nil},
		{Platform{OS: windows, Arch: amd64}, v3, "Godot_v3.6-beta1_win64.exe", nil},
		{Platform{OS: windows, Arch: arm64}, v3, "", ErrUnsupportedArch},
		{Platform{OS: windows, Arch: universal}, v3, "", ErrUnsupportedArch},

		// v4.0
		{Platform{OS: windows, Arch: i386}, v4, "Godot_v4.0-stable_win32.exe", nil},
		{Platform{OS: windows, Arch: amd64}, v4, "Godot_v4.0-stable_win64.exe", nil},
		{Platform{OS: windows, Arch: arm64}, v4, "", ErrUnsupportedArch},
		{Platform{OS: windows, Arch: universal}, v4, "", ErrUnsupportedArch},

		// v4.0-stable_mono
		{Platform{OS: windows, Arch: i386}, v4Mono, "Godot_v4.0-stable_mono_win32.exe", nil},
		{Platform{OS: windows, Arch: amd64}, v4Mono, "Godot_v4.0-stable_mono_win64.exe", nil},
		{Platform{OS: windows, Arch: arm64}, v4Mono, "", ErrUnsupportedArch},
		{Platform{OS: windows, Arch: universal}, v4Mono, "", ErrUnsupportedArch},

		// v5.0-rc4
		{Platform{OS: windows, Arch: i386}, v5, "Godot_v5.0-rc4_win32.exe", nil},
		{Platform{OS: windows, Arch: amd64}, v5, "Godot_v5.0-rc4_win64.exe", nil},
		{Platform{OS: windows, Arch: arm64}, v5, "", ErrUnsupportedArch},
		{Platform{OS: windows, Arch: universal}, v5, "", ErrUnsupportedArch},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%v-%s", tc.p, tc.v), func(t *testing.T) {
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
	var (
		v1     = Version{major: 1}
		v2     = Version{major: 2, label: "beta10"}
		v3     = Version{3, 0, 4, "alpha1"}
		v4     = Version{4, 0, 11, "dev.20230101"}
		v4Mono = Version{4, 0, 0, "stable_mono"}
	)

	tests := []struct {
		s    string
		want Executable
		err  error
	}{
		// Invalid inputs
		{s: "", want: Executable{}, err: ErrMissingName},
		{s: "Godot-v1.0-stable-x11.32", want: Executable{}, err: ErrInvalidName},
		{s: "Godot-v1.0_stable-x11.32", want: Executable{}, err: ErrInvalidName},
		{s: "Godot_invalid_x11.32", want: Executable{}, err: ErrInvalidVersion},
		{s: "Godot_v1.0-stable_invalid", want: Executable{}, err: ErrUnrecognizedPlatform},

		// Valid inputs
		// Linux
		{s: "Godot_v1.0-stable_x11.32", want: Executable{Platform{i386, linux}, v1}, err: nil},
		{s: "Godot_v2.0-beta10_x11.64", want: Executable{Platform{amd64, linux}, v2}, err: nil},
		{s: "Godot_v3.0.4-alpha1_x11.32", want: Executable{Platform{i386, linux}, v3}, err: nil},
		{s: "Godot_v4.0.11-dev.20230101_x11.64", want: Executable{Platform{amd64, linux}, v4}, err: nil},
		{s: "Godot_v4.0-stable_mono_linux_x86_64", want: Executable{Platform{amd64, linux}, v4Mono}, err: nil},

		// Darwin
		{s: "Godot_v1.0-stable_osx.fat", want: Executable{Platform{universal, macOS}, v1}, err: nil},
		{s: "Godot_v2.0-beta10_osx.64", want: Executable{Platform{amd64, macOS}, v2}, err: nil},
		{s: "Godot_v3.0.4-alpha1_osx.universal", want: Executable{Platform{universal, macOS}, v3}, err: nil},
		{s: "Godot_v4.0.11-dev.20230101_macos.universal", want: Executable{Platform{universal, macOS}, v4}, err: nil},
		{s: "Godot_v4.0-stable_mono_macos.universal", want: Executable{Platform{universal, macOS}, v4Mono}, err: nil},

		// Windows
		{s: "Godot_v1.0-stable_win32", want: Executable{Platform{i386, windows}, v1}, err: nil},
		{s: "Godot_v2.0-beta10_win64", want: Executable{Platform{amd64, windows}, v2}, err: nil},
		{s: "Godot_v3.0.4-alpha1_win32", want: Executable{Platform{i386, windows}, v3}, err: nil},
		{s: "Godot_v4.0.11-dev.20230101_win64", want: Executable{Platform{amd64, windows}, v4}, err: nil},
		{s: "Godot_v4.0-stable_mono_win64", want: Executable{Platform{amd64, windows}, v4Mono}, err: nil},
	}

	for _, tc := range tests {
		t.Run(tc.s, func(t *testing.T) {
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
