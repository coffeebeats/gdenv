package godot

import (
	"errors"
	"fmt"
	"testing"
)

/* ------------------------------ Test: ParseOS ----------------------------- */

func TestParseOS(t *testing.T) {
	tests := []struct {
		s    string
		want OS
		err  error
	}{
		// Invalid inputs
		{s: "", err: ErrMissingOSInput},
		{s: "abc", err: ErrUnrecognizedOSInput},
		{s: "linux-", err: ErrUnrecognizedOSInput},
		{s: "mac.os", err: ErrUnrecognizedOSInput},
		{s: "win32", err: ErrUnrecognizedOSInput},

		// Valid inputs (Go-defined)
		{s: "linux", want: Linux},

		{s: "darwin", want: MacOS},
		{s: "macos", want: MacOS},
		{s: "osx", want: MacOS},

		{s: "win", want: Windows},
		{s: "windows", want: Windows},

		// Valid inputs (user-supplied)
		{s: "LINUX", want: Linux},
		{s: " LINUX\n", want: Linux},
		{s: "\tOSX ", want: MacOS},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got, err := ParseOS(tc.s)

			if !errors.Is(err, tc.err) {
				t.Fatalf("err: got %#v, want %#v", err, tc.err)
			}
			if got != tc.want {
				t.Fatalf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}

/* ----------------------------- Test: ParseArch ---------------------------- */

func TestParseArch(t *testing.T) {
	tests := []struct {
		s    string
		want Arch
		err  error
	}{
		// Invalid inputs
		{s: "", err: ErrMissingArchInput},
		{s: "abc", err: ErrUnrecognizedArchInput},

		// Valid inputs (Go-defined)
		{s: "amd64", want: Amd64},
		{s: "x86_64", want: Amd64},
		{s: "x86-64", want: Amd64},

		{s: "arm64", want: Arm64},
		{s: "arm64be", want: Arm64},

		{s: "386", want: I386},
		{s: "i386", want: I386},
		{s: "x86", want: I386},

		{s: "fat", want: Universal},
		{s: "universal", want: Universal},

		// Valid inputs (user-supplied)
		{s: "AMD64", want: Amd64},
		{s: " X86_64\n", want: Amd64},
		{s: "\tuniversal ", want: Universal},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got, err := ParseArch(tc.s)

			if !errors.Is(err, tc.err) {
				t.Fatalf("err: got %#v, want %#v", err, tc.err)
			}
			if got != tc.want {
				t.Fatalf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}

/* -------------------------- Test: FormatPlatform -------------------------- */

func TestFormatPlatform(t *testing.T) {
	tests := []struct {
		p    Platform
		v    Version
		want string
		err  error
	}{
		// Invalid inputs
		{Platform{}, Version{}, "", ErrMissingOSInput},
		{Platform{os: Linux}, Version{}, "", ErrMissingArchInput},
		{Platform{os: Linux, arch: Amd64}, Version{}, "", ErrUnsupportedVersion},

		{Platform{os: Linux, arch: Amd64}, Version{major: 2}, "", ErrUnsupportedVersion},
		{Platform{os: MacOS, arch: Amd64}, Version{major: 2}, "", ErrUnsupportedVersion},
		{Platform{os: Windows, arch: Amd64}, Version{major: 2}, "", ErrUnsupportedVersion},

		// Valid inputs - Linux

		// v3.*
		{Platform{os: Linux, arch: I386}, Version{major: 3}, "x11.32", nil},
		{Platform{os: Linux, arch: Amd64}, Version{major: 3}, "x11.64", nil},
		{Platform{os: Linux, arch: Arm64}, Version{major: 3}, "", ErrUnsupportedArchInput},
		{Platform{os: Linux, arch: Universal}, Version{major: 3}, "", ErrUnsupportedArchInput},

		// v4.0+
		{Platform{os: Linux, arch: I386}, Version{major: 4}, "linux.x86_32", nil},
		{Platform{os: Linux, arch: Amd64}, Version{major: 4}, "linux.x86_64", nil},
		{Platform{os: Linux, arch: Arm64}, Version{major: 4}, "", ErrUnsupportedArchInput},
		{Platform{os: Linux, arch: Universal}, Version{major: 4}, "", ErrUnsupportedArchInput},

		// Valid inputs - MacOS

		// v3.0 - v3.0.6
		{Platform{os: MacOS, arch: I386}, Version{major: 3}, "osx.fat", nil},
		{Platform{os: MacOS, arch: Amd64}, Version{major: 3}, "osx.fat", nil},
		{Platform{os: MacOS, arch: Arm64}, Version{major: 3}, "", ErrUnsupportedArchInput},
		{Platform{os: MacOS, arch: Universal}, Version{major: 3}, "", ErrUnsupportedArchInput},

		{Platform{os: MacOS, arch: I386}, Version{major: 3, patch: 6}, "osx.fat", nil},
		{Platform{os: MacOS, arch: Amd64}, Version{major: 3, patch: 6}, "osx.fat", nil},
		{Platform{os: MacOS, arch: Arm64}, Version{major: 3, patch: 6}, "", ErrUnsupportedArchInput},
		{Platform{os: MacOS, arch: Universal}, Version{major: 3, patch: 6}, "", ErrUnsupportedArchInput},

		// v3.1 - v3.2.4-beta2
		{Platform{os: MacOS, arch: Amd64}, Version{major: 3, minor: 1}, "osx.64", nil},
		{Platform{os: MacOS, arch: I386}, Version{major: 3, minor: 1}, "", ErrUnsupportedArchInput},
		{Platform{os: MacOS, arch: Arm64}, Version{major: 3, minor: 1}, "", ErrUnsupportedArchInput},
		{Platform{os: MacOS, arch: Universal}, Version{major: 3, minor: 1}, "", ErrUnsupportedArchInput},

		{Platform{os: MacOS, arch: Amd64}, Version{3, 2, 4, "beta2"}, "osx.64", nil},
		{Platform{os: MacOS, arch: I386}, Version{3, 2, 4, "beta2"}, "", ErrUnsupportedArchInput},
		{Platform{os: MacOS, arch: Arm64}, Version{3, 2, 4, "beta2"}, "", ErrUnsupportedArchInput},
		{Platform{os: MacOS, arch: Universal}, Version{3, 2, 4, "beta2"}, "", ErrUnsupportedArchInput},

		// v3.2.4-beta3 - v4.0-alpha12
		{Platform{os: MacOS, arch: Amd64}, Version{3, 2, 4, "beta3"}, "osx.universal", nil},
		{Platform{os: MacOS, arch: Arm64}, Version{3, 2, 4, "beta3"}, "osx.universal", nil},
		{Platform{os: MacOS, arch: I386}, Version{3, 2, 4, "beta3"}, "", ErrUnsupportedArchInput},
		{Platform{os: MacOS, arch: Universal}, Version{3, 2, 4, "beta3"}, "", ErrUnsupportedArchInput},

		{Platform{os: MacOS, arch: Amd64}, Version{3, 2, 4, "rc1"}, "osx.universal", nil},
		{Platform{os: MacOS, arch: Arm64}, Version{3, 2, 4, "rc1"}, "osx.universal", nil},
		{Platform{os: MacOS, arch: I386}, Version{3, 2, 4, "rc1"}, "", ErrUnsupportedArchInput},
		{Platform{os: MacOS, arch: Universal}, Version{3, 2, 4, "rc1"}, "", ErrUnsupportedArchInput},

		{Platform{os: MacOS, arch: Amd64}, Version{3, 2, 4, "stable"}, "osx.universal", nil},
		{Platform{os: MacOS, arch: Arm64}, Version{3, 2, 4, "stable"}, "osx.universal", nil},
		{Platform{os: MacOS, arch: I386}, Version{3, 2, 4, "stable"}, "", ErrUnsupportedArchInput},
		{Platform{os: MacOS, arch: Universal}, Version{3, 2, 4, "stable"}, "", ErrUnsupportedArchInput},

		{Platform{os: MacOS, arch: Amd64}, Version{4, 0, 0, "dev.20210727"}, "osx.universal", nil},
		{Platform{os: MacOS, arch: Arm64}, Version{4, 0, 0, "dev.20210727"}, "osx.universal", nil},
		{Platform{os: MacOS, arch: I386}, Version{4, 0, 0, "dev.20210727"}, "", ErrUnsupportedArchInput},
		{Platform{os: MacOS, arch: Universal}, Version{4, 0, 0, "dev.20210727"}, "", ErrUnsupportedArchInput},

		{Platform{os: MacOS, arch: Amd64}, Version{4, 0, 0, "alpha1"}, "osx.universal", nil},
		{Platform{os: MacOS, arch: Arm64}, Version{4, 0, 0, "alpha1"}, "osx.universal", nil},
		{Platform{os: MacOS, arch: I386}, Version{4, 0, 0, "alpha1"}, "", ErrUnsupportedArchInput},
		{Platform{os: MacOS, arch: Universal}, Version{4, 0, 0, "alpha1"}, "", ErrUnsupportedArchInput},

		{Platform{os: MacOS, arch: Amd64}, Version{4, 0, 0, "alpha12"}, "osx.universal", nil},
		{Platform{os: MacOS, arch: Arm64}, Version{4, 0, 0, "alpha12"}, "osx.universal", nil},
		{Platform{os: MacOS, arch: I386}, Version{4, 0, 0, "alpha12"}, "", ErrUnsupportedArchInput},
		{Platform{os: MacOS, arch: Universal}, Version{4, 0, 0, "alpha12"}, "", ErrUnsupportedArchInput},

		// v4.0-alpha13+
		{Platform{os: MacOS, arch: Amd64}, Version{4, 0, 0, "alpha13"}, "macos.universal", nil},
		{Platform{os: MacOS, arch: Arm64}, Version{4, 0, 0, "alpha13"}, "macos.universal", nil},
		{Platform{os: MacOS, arch: I386}, Version{4, 0, 0, "alpha13"}, "", ErrUnsupportedArchInput},
		{Platform{os: MacOS, arch: Universal}, Version{4, 0, 0, "alpha13"}, "", ErrUnsupportedArchInput},

		{Platform{os: MacOS, arch: Amd64}, Version{4, 0, 0, "beta1"}, "macos.universal", nil},
		{Platform{os: MacOS, arch: Arm64}, Version{4, 0, 0, "beta1"}, "macos.universal", nil},
		{Platform{os: MacOS, arch: I386}, Version{4, 0, 0, "beta1"}, "", ErrUnsupportedArchInput},
		{Platform{os: MacOS, arch: Universal}, Version{4, 0, 0, "beta1"}, "", ErrUnsupportedArchInput},

		{Platform{os: MacOS, arch: Amd64}, Version{4, 0, 0, "rc1"}, "macos.universal", nil},
		{Platform{os: MacOS, arch: Arm64}, Version{4, 0, 0, "rc1"}, "macos.universal", nil},
		{Platform{os: MacOS, arch: I386}, Version{4, 0, 0, "rc1"}, "", ErrUnsupportedArchInput},
		{Platform{os: MacOS, arch: Universal}, Version{4, 0, 0, "rc1"}, "", ErrUnsupportedArchInput},

		{Platform{os: MacOS, arch: Amd64}, Version{4, 0, 0, "stable"}, "macos.universal", nil},
		{Platform{os: MacOS, arch: Arm64}, Version{4, 0, 0, "stable"}, "macos.universal", nil},
		{Platform{os: MacOS, arch: I386}, Version{4, 0, 0, "stable"}, "", ErrUnsupportedArchInput},
		{Platform{os: MacOS, arch: Universal}, Version{4, 0, 0, "stable"}, "", ErrUnsupportedArchInput},

		// Valid inputs - Windows

		// v3.*
		{Platform{os: Windows, arch: I386}, Version{major: 3}, "win32", nil},
		{Platform{os: Windows, arch: Amd64}, Version{major: 3}, "win64", nil},
		{Platform{os: Windows, arch: Arm64}, Version{major: 3}, "", ErrUnsupportedArchInput},
		{Platform{os: Windows, arch: Universal}, Version{major: 3}, "", ErrUnsupportedArchInput},

		// v4.0+
		{Platform{os: Windows, arch: I386}, Version{major: 4}, "win32", nil},
		{Platform{os: Windows, arch: Amd64}, Version{major: 4}, "win64", nil},
		{Platform{os: Windows, arch: Arm64}, Version{major: 4}, "", ErrUnsupportedArchInput},
		{Platform{os: Windows, arch: Universal}, Version{major: 4}, "", ErrUnsupportedArchInput},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%v %s", tc.p, tc.v), func(t *testing.T) {
			got, err := FormatPlatform(tc.p, tc.v)

			if !errors.Is(err, tc.err) {
				t.Fatalf("err: got %#v, want %#v", err, tc.err)
			}
			if got != tc.want {
				t.Fatalf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}
