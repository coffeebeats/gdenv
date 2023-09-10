package godot

import (
	"errors"
	"fmt"
	"testing"

	"github.com/coffeebeats/gdenv/internal/version"
)

/* -------------------------- Test: Executable.Name ------------------------- */

func TestExecutableName(t *testing.T) {
	var (
		v3     = version.MustParse("3.6-beta1")
		v4     = version.MustParse("4")
		v4Mono = version.MustParse("4.0-stable_mono")
		v5     = version.MustParse("5.0-rc4")
	)

	tests := []struct {
		p    Platform
		v    version.Version
		want string
		err  error
	}{
		// Invalid inputs
		{Platform{}, version.Version{}, "", ErrMissingOS},
		{Platform{OS: Linux}, version.Version{}, "", ErrMissingArch},
		{Platform{OS: Linux, Arch: Amd64}, version.Version{}, "", version.ErrUnsupported},

		{Platform{OS: Linux, Arch: Amd64}, version.MustParse("2"), "", version.ErrUnsupported},
		{Platform{OS: MacOS, Arch: Amd64}, version.MustParse("2"), "", version.ErrUnsupported},
		{Platform{OS: Windows, Arch: Amd64}, version.MustParse("2"), "", version.ErrUnsupported},

		// Valid inputs - Linux

		// v3.6-beta1
		{Platform{OS: Linux, Arch: I386}, v3, "Godot_v3.6-beta1_x11.32", nil},
		{Platform{OS: Linux, Arch: Amd64}, v3, "Godot_v3.6-beta1_x11.64", nil},
		{Platform{OS: Linux, Arch: Arm64}, v3, "", ErrUnsupportedArch},
		{Platform{OS: Linux, Arch: Universal}, v3, "", ErrUnsupportedArch},

		// v4.0
		{Platform{OS: Linux, Arch: I386}, v4, "Godot_v4.0-stable_linux.x86_32", nil},
		{Platform{OS: Linux, Arch: Amd64}, v4, "Godot_v4.0-stable_linux.x86_64", nil},
		{Platform{OS: Linux, Arch: Arm64}, v4, "", ErrUnsupportedArch},
		{Platform{OS: Linux, Arch: Universal}, v4, "", ErrUnsupportedArch},

		// v4.0-stable_mono
		{Platform{OS: Linux, Arch: I386}, v4Mono, "Godot_v4.0-stable_mono_linux_x86_32", nil},
		{Platform{OS: Linux, Arch: Amd64}, v4Mono, "Godot_v4.0-stable_mono_linux_x86_64", nil},
		{Platform{OS: Linux, Arch: Arm64}, v4Mono, "", ErrUnsupportedArch},
		{Platform{OS: Linux, Arch: Universal}, v4Mono, "", ErrUnsupportedArch},

		// v5.0-rc4
		{Platform{OS: Linux, Arch: I386}, v5, "Godot_v5.0-rc4_linux.x86_32", nil},
		{Platform{OS: Linux, Arch: Amd64}, v5, "Godot_v5.0-rc4_linux.x86_64", nil},
		{Platform{OS: Linux, Arch: Arm64}, v5, "", ErrUnsupportedArch},
		{Platform{OS: Linux, Arch: Universal}, v5, "", ErrUnsupportedArch},

		// Valid inputs - MacOS

		// v3.6-beta1
		{Platform{OS: MacOS, Arch: Amd64}, v3, "Godot_v3.6-beta1_osx.universal", nil},
		{Platform{OS: MacOS, Arch: Arm64}, v3, "Godot_v3.6-beta1_osx.universal", nil},
		{Platform{OS: MacOS, Arch: I386}, v3, "", ErrUnsupportedArch},
		{Platform{OS: MacOS, Arch: Universal}, v3, "", ErrUnsupportedArch},

		// v4.0
		{Platform{OS: MacOS, Arch: Amd64}, v4, "Godot_v4.0-stable_macos.universal", nil},
		{Platform{OS: MacOS, Arch: Arm64}, v4, "Godot_v4.0-stable_macos.universal", nil},
		{Platform{OS: MacOS, Arch: I386}, v4, "", ErrUnsupportedArch},
		{Platform{OS: MacOS, Arch: Universal}, v4, "", ErrUnsupportedArch},

		// v4.0-stable_mono
		{Platform{OS: MacOS, Arch: Amd64}, v4Mono, "Godot_v4.0-stable_mono_macos.universal", nil},
		{Platform{OS: MacOS, Arch: Arm64}, v4Mono, "Godot_v4.0-stable_mono_macos.universal", nil},
		{Platform{OS: MacOS, Arch: I386}, v4Mono, "", ErrUnsupportedArch},
		{Platform{OS: MacOS, Arch: Universal}, v4Mono, "", ErrUnsupportedArch},

		// v5.0-rc4
		{Platform{OS: MacOS, Arch: Amd64}, v5, "Godot_v5.0-rc4_macos.universal", nil},
		{Platform{OS: MacOS, Arch: Arm64}, v5, "Godot_v5.0-rc4_macos.universal", nil},
		{Platform{OS: MacOS, Arch: I386}, v5, "", ErrUnsupportedArch},
		{Platform{OS: MacOS, Arch: Universal}, v5, "", ErrUnsupportedArch},

		// Valid inputs - Windows

		// v3.6-beta1
		{Platform{OS: Windows, Arch: I386}, v3, "Godot_v3.6-beta1_win32.exe", nil},
		{Platform{OS: Windows, Arch: Amd64}, v3, "Godot_v3.6-beta1_win64.exe", nil},
		{Platform{OS: Windows, Arch: Arm64}, v3, "", ErrUnsupportedArch},
		{Platform{OS: Windows, Arch: Universal}, v3, "", ErrUnsupportedArch},

		// v4.0
		{Platform{OS: Windows, Arch: I386}, v4, "Godot_v4.0-stable_win32.exe", nil},
		{Platform{OS: Windows, Arch: Amd64}, v4, "Godot_v4.0-stable_win64.exe", nil},
		{Platform{OS: Windows, Arch: Arm64}, v4, "", ErrUnsupportedArch},
		{Platform{OS: Windows, Arch: Universal}, v4, "", ErrUnsupportedArch},

		// v4.0-stable_mono
		{Platform{OS: Windows, Arch: I386}, v4Mono, "Godot_v4.0-stable_mono_win32.exe", nil},
		{Platform{OS: Windows, Arch: Amd64}, v4Mono, "Godot_v4.0-stable_mono_win64.exe", nil},
		{Platform{OS: Windows, Arch: Arm64}, v4Mono, "", ErrUnsupportedArch},
		{Platform{OS: Windows, Arch: Universal}, v4Mono, "", ErrUnsupportedArch},

		// v5.0-rc4
		{Platform{OS: Windows, Arch: I386}, v5, "Godot_v5.0-rc4_win32.exe", nil},
		{Platform{OS: Windows, Arch: Amd64}, v5, "Godot_v5.0-rc4_win64.exe", nil},
		{Platform{OS: Windows, Arch: Arm64}, v5, "", ErrUnsupportedArch},
		{Platform{OS: Windows, Arch: Universal}, v5, "", ErrUnsupportedArch},
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
		v1     = version.MustParse("1")
		v2     = version.MustParse("2.0-beta10")
		v3     = version.MustParse("3.0.4-alpha1")
		v4     = version.MustParse("4.0.11-dev.20230101")
		v4Mono = version.MustParse("4.0-stable_mono")
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
		{s: "Godot_invalid_x11.32", want: Executable{}, err: version.ErrInvalid},
		{s: "Godot_v1.0-stable_invalid", want: Executable{}, err: ErrUnrecognizedPlatform},

		// Valid inputs
		// Linux
		{s: "Godot_v1.0-stable_x11.32", want: Executable{Platform{I386, Linux}, v1}, err: nil},
		{s: "Godot_v2.0-beta10_x11.64", want: Executable{Platform{Amd64, Linux}, v2}, err: nil},
		{s: "Godot_v3.0.4-alpha1_x11.32", want: Executable{Platform{I386, Linux}, v3}, err: nil},
		{s: "Godot_v4.0.11-dev.20230101_x11.64", want: Executable{Platform{Amd64, Linux}, v4}, err: nil},
		{s: "Godot_v4.0-stable_mono_linux_x86_64", want: Executable{Platform{Amd64, Linux}, v4Mono}, err: nil},

		// Darwin
		{s: "Godot_v1.0-stable_osx.fat", want: Executable{Platform{Universal, MacOS}, v1}, err: nil},
		{s: "Godot_v2.0-beta10_osx.64", want: Executable{Platform{Amd64, MacOS}, v2}, err: nil},
		{s: "Godot_v3.0.4-alpha1_osx.universal", want: Executable{Platform{Universal, MacOS}, v3}, err: nil},
		{s: "Godot_v4.0.11-dev.20230101_macos.universal", want: Executable{Platform{Universal, MacOS}, v4}, err: nil},
		{s: "Godot_v4.0-stable_mono_macos.universal", want: Executable{Platform{Universal, MacOS}, v4Mono}, err: nil},

		// Windows
		{s: "Godot_v1.0-stable_win32", want: Executable{Platform{I386, Windows}, v1}, err: nil},
		{s: "Godot_v2.0-beta10_win64", want: Executable{Platform{Amd64, Windows}, v2}, err: nil},
		{s: "Godot_v3.0.4-alpha1_win32", want: Executable{Platform{I386, Windows}, v3}, err: nil},
		{s: "Godot_v4.0.11-dev.20230101_win64", want: Executable{Platform{Amd64, Windows}, v4}, err: nil},
		{s: "Godot_v4.0-stable_mono_win64", want: Executable{Platform{Amd64, Windows}, v4Mono}, err: nil},
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
