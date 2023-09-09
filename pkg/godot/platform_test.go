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
		{s: "", err: ErrMissingOS},
		{s: "abc", err: ErrUnrecognizedOS},
		{s: "linux-", err: ErrUnrecognizedOS},
		{s: "mac.os", err: ErrUnrecognizedOS},
		{s: "win32", err: ErrUnrecognizedOS},

		// Valid inputs (Go-defined)
		{s: "linux", want: linux},

		{s: "darwin", want: macOS},
		{s: "macos", want: macOS},
		{s: "osx", want: macOS},

		{s: "win", want: windows},
		{s: "windows", want: windows},

		// Valid inputs (user-supplied)
		{s: "LINUX", want: linux},
		{s: " LINUX\n", want: linux},
		{s: "\tOSX ", want: macOS},
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
		{s: "", err: ErrMissingArch},
		{s: "abc", err: ErrUnrecognizedArch},

		// Valid inputs (Go-defined)
		{s: "amd64", want: amd64},
		{s: "x86_64", want: amd64},
		{s: "x86-64", want: amd64},

		{s: "arm64", want: arm64},
		{s: "arm64be", want: arm64},

		{s: "386", want: i386},
		{s: "i386", want: i386},
		{s: "x86", want: i386},

		{s: "fat", want: universal},
		{s: "universal", want: universal},

		// Valid inputs (user-supplied)
		{s: "AMD64", want: amd64},
		{s: " X86_64\n", want: amd64},
		{s: "\tuniversal ", want: universal},
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

/* --------------------------- Test: ParsePlatform -------------------------- */

func TestParsePlatform(t *testing.T) {
	tests := []struct {
		s    string
		want Platform
		err  error
	}{
		// Invalid inputs
		{s: "", err: ErrMissingPlatform},
		{s: "abc", err: ErrUnrecognizedPlatform},

		// Valid inputs (Go-defined)
		// Linux
		{s: "x11.32", want: Platform{i386, linux}, err: nil},
		{s: "x11.64", want: Platform{amd64, linux}, err: nil},
		{s: "linux.x86_32", want: Platform{i386, linux}, err: nil},
		{s: "linux.x86_64", want: Platform{amd64, linux}, err: nil},
		{s: "linux_x86_32", want: Platform{i386, linux}, err: nil},
		{s: "linux_x86_64", want: Platform{amd64, linux}, err: nil},

		// MacOS
		{s: "osx.64", want: Platform{amd64, macOS}, err: nil},
		{s: "macos.universal", want: Platform{universal, macOS}, err: nil},
		{s: "osx.fat", want: Platform{universal, macOS}, err: nil},
		{s: "osx.universal", want: Platform{universal, macOS}, err: nil},

		// Windows
		{s: "win32", want: Platform{i386, windows}, err: nil},
		{s: "win64", want: Platform{amd64, windows}, err: nil},

		// Valid inputs (user-supplied)
		{s: "WIN64", want: Platform{amd64, windows}, err: nil},
		{s: " MACOS.UNIVERSAL\n", want: Platform{universal, macOS}},
		{s: "\tlinux.x86_64 ", want: Platform{amd64, linux}, err: nil},
	}

	for _, tc := range tests {
		t.Run(tc.s, func(t *testing.T) {
			got, err := ParsePlatform(tc.s)

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
		{Platform{}, Version{}, "", ErrMissingOS},
		{Platform{os: linux}, Version{}, "", ErrMissingArch},
		{Platform{os: linux, arch: amd64}, Version{}, "", ErrUnsupportedVersion},

		{Platform{os: linux, arch: amd64}, Version{major: 2}, "", ErrUnsupportedVersion},
		{Platform{os: macOS, arch: amd64}, Version{major: 2}, "", ErrUnsupportedVersion},
		{Platform{os: windows, arch: amd64}, Version{major: 2}, "", ErrUnsupportedVersion},

		// Valid inputs - linux

		// v3.*
		{Platform{os: linux, arch: i386}, Version{major: 3}, "x11.32", nil},
		{Platform{os: linux, arch: amd64}, Version{major: 3}, "x11.64", nil},
		{Platform{os: linux, arch: i386}, Version{major: 3, label: LabelMono}, "x11_32", nil},
		{Platform{os: linux, arch: amd64}, Version{major: 3, label: LabelMono}, "x11_64", nil},
		{Platform{os: linux, arch: arm64}, Version{major: 3}, "", ErrUnsupportedArch},
		{Platform{os: linux, arch: universal}, Version{major: 3}, "", ErrUnsupportedArch},

		// v4.0-dev.* - v4.0-alpha14
		{Platform{os: linux, arch: i386}, Version{major: 4, label: "dev.20220118"}, "linux.32", nil},
		{Platform{os: linux, arch: amd64}, Version{major: 4, label: "dev.20220118"}, "linux.64", nil},
		{Platform{os: linux, arch: arm64}, Version{major: 4, label: "dev.20220118"}, "", ErrUnsupportedArch},
		{Platform{os: linux, arch: universal}, Version{major: 4, label: "dev.20220118"}, "", ErrUnsupportedArch},

		{Platform{os: linux, arch: i386}, Version{major: 4, label: "alpha14"}, "linux.32", nil},
		{Platform{os: linux, arch: amd64}, Version{major: 4, label: "alpha14"}, "linux.64", nil},
		{Platform{os: linux, arch: arm64}, Version{major: 4, label: "alpha14"}, "", ErrUnsupportedArch},
		{Platform{os: linux, arch: universal}, Version{major: 4, label: "alpha14"}, "", ErrUnsupportedArch},

		// v4.0-alpha15+
		{Platform{os: linux, arch: i386}, Version{major: 4, label: "alpha15"}, "linux.x86_32", nil},
		{Platform{os: linux, arch: amd64}, Version{major: 4, label: "alpha15"}, "linux.x86_64", nil},
		{Platform{os: linux, arch: arm64}, Version{major: 4, label: "alpha15"}, "", ErrUnsupportedArch},
		{Platform{os: linux, arch: universal}, Version{major: 4, label: "alpha15"}, "", ErrUnsupportedArch},

		{Platform{os: linux, arch: i386}, Version{major: 4}, "linux.x86_32", nil},
		{Platform{os: linux, arch: amd64}, Version{major: 4}, "linux.x86_64", nil},
		{Platform{os: linux, arch: i386}, Version{major: 4, label: LabelMono}, "linux_x86_32", nil},
		{Platform{os: linux, arch: amd64}, Version{major: 4, label: LabelMono}, "linux_x86_64", nil},
		{Platform{os: linux, arch: arm64}, Version{major: 4}, "", ErrUnsupportedArch},
		{Platform{os: linux, arch: universal}, Version{major: 4}, "", ErrUnsupportedArch},

		// Valid inputs - MacOS

		// v3.0 - v3.0.6
		{Platform{os: macOS, arch: i386}, Version{major: 3}, "osx.fat", nil},
		{Platform{os: macOS, arch: amd64}, Version{major: 3}, "osx.fat", nil},
		{Platform{os: macOS, arch: arm64}, Version{major: 3}, "", ErrUnsupportedArch},
		{Platform{os: macOS, arch: universal}, Version{major: 3}, "", ErrUnsupportedArch},

		{Platform{os: macOS, arch: i386}, Version{major: 3, patch: 6}, "osx.fat", nil},
		{Platform{os: macOS, arch: amd64}, Version{major: 3, patch: 6}, "osx.fat", nil},
		{Platform{os: macOS, arch: arm64}, Version{major: 3, patch: 6}, "", ErrUnsupportedArch},
		{Platform{os: macOS, arch: universal}, Version{major: 3, patch: 6}, "", ErrUnsupportedArch},

		// v3.1 - v3.2.4-beta2
		{Platform{os: macOS, arch: amd64}, Version{major: 3, minor: 1}, "osx.64", nil},
		{Platform{os: macOS, arch: i386}, Version{major: 3, minor: 1}, "", ErrUnsupportedArch},
		{Platform{os: macOS, arch: arm64}, Version{major: 3, minor: 1}, "", ErrUnsupportedArch},
		{Platform{os: macOS, arch: universal}, Version{major: 3, minor: 1}, "", ErrUnsupportedArch},

		{Platform{os: macOS, arch: amd64}, Version{3, 2, 4, "beta2"}, "osx.64", nil},
		{Platform{os: macOS, arch: i386}, Version{3, 2, 4, "beta2"}, "", ErrUnsupportedArch},
		{Platform{os: macOS, arch: arm64}, Version{3, 2, 4, "beta2"}, "", ErrUnsupportedArch},
		{Platform{os: macOS, arch: universal}, Version{3, 2, 4, "beta2"}, "", ErrUnsupportedArch},

		// v3.2.4-beta3 - v4.0-alpha12
		{Platform{os: macOS, arch: amd64}, Version{3, 2, 4, "beta3"}, "osx.universal", nil},
		{Platform{os: macOS, arch: arm64}, Version{3, 2, 4, "beta3"}, "osx.universal", nil},
		{Platform{os: macOS, arch: i386}, Version{3, 2, 4, "beta3"}, "", ErrUnsupportedArch},
		{Platform{os: macOS, arch: universal}, Version{3, 2, 4, "beta3"}, "", ErrUnsupportedArch},

		{Platform{os: macOS, arch: amd64}, Version{3, 2, 4, "rc1"}, "osx.universal", nil},
		{Platform{os: macOS, arch: arm64}, Version{3, 2, 4, "rc1"}, "osx.universal", nil},
		{Platform{os: macOS, arch: i386}, Version{3, 2, 4, "rc1"}, "", ErrUnsupportedArch},
		{Platform{os: macOS, arch: universal}, Version{3, 2, 4, "rc1"}, "", ErrUnsupportedArch},

		{Platform{os: macOS, arch: amd64}, Version{3, 2, 4, "stable"}, "osx.universal", nil},
		{Platform{os: macOS, arch: arm64}, Version{3, 2, 4, "stable"}, "osx.universal", nil},
		{Platform{os: macOS, arch: i386}, Version{3, 2, 4, "stable"}, "", ErrUnsupportedArch},
		{Platform{os: macOS, arch: universal}, Version{3, 2, 4, "stable"}, "", ErrUnsupportedArch},

		{Platform{os: macOS, arch: amd64}, Version{4, 0, 0, "dev.20210727"}, "osx.universal", nil},
		{Platform{os: macOS, arch: arm64}, Version{4, 0, 0, "dev.20210727"}, "osx.universal", nil},
		{Platform{os: macOS, arch: i386}, Version{4, 0, 0, "dev.20210727"}, "", ErrUnsupportedArch},
		{Platform{os: macOS, arch: universal}, Version{4, 0, 0, "dev.20210727"}, "", ErrUnsupportedArch},

		{Platform{os: macOS, arch: amd64}, Version{4, 0, 0, "alpha1"}, "osx.universal", nil},
		{Platform{os: macOS, arch: arm64}, Version{4, 0, 0, "alpha1"}, "osx.universal", nil},
		{Platform{os: macOS, arch: i386}, Version{4, 0, 0, "alpha1"}, "", ErrUnsupportedArch},
		{Platform{os: macOS, arch: universal}, Version{4, 0, 0, "alpha1"}, "", ErrUnsupportedArch},

		{Platform{os: macOS, arch: amd64}, Version{4, 0, 0, "alpha12"}, "osx.universal", nil},
		{Platform{os: macOS, arch: arm64}, Version{4, 0, 0, "alpha12"}, "osx.universal", nil},
		{Platform{os: macOS, arch: i386}, Version{4, 0, 0, "alpha12"}, "", ErrUnsupportedArch},
		{Platform{os: macOS, arch: universal}, Version{4, 0, 0, "alpha12"}, "", ErrUnsupportedArch},

		// v4.0-alpha13+
		{Platform{os: macOS, arch: amd64}, Version{4, 0, 0, "alpha13"}, "macos.universal", nil},
		{Platform{os: macOS, arch: arm64}, Version{4, 0, 0, "alpha13"}, "macos.universal", nil},
		{Platform{os: macOS, arch: i386}, Version{4, 0, 0, "alpha13"}, "", ErrUnsupportedArch},
		{Platform{os: macOS, arch: universal}, Version{4, 0, 0, "alpha13"}, "", ErrUnsupportedArch},

		{Platform{os: macOS, arch: amd64}, Version{4, 0, 0, "beta1"}, "macos.universal", nil},
		{Platform{os: macOS, arch: arm64}, Version{4, 0, 0, "beta1"}, "macos.universal", nil},
		{Platform{os: macOS, arch: i386}, Version{4, 0, 0, "beta1"}, "", ErrUnsupportedArch},
		{Platform{os: macOS, arch: universal}, Version{4, 0, 0, "beta1"}, "", ErrUnsupportedArch},

		{Platform{os: macOS, arch: amd64}, Version{4, 0, 0, "rc1"}, "macos.universal", nil},
		{Platform{os: macOS, arch: arm64}, Version{4, 0, 0, "rc1"}, "macos.universal", nil},
		{Platform{os: macOS, arch: i386}, Version{4, 0, 0, "rc1"}, "", ErrUnsupportedArch},
		{Platform{os: macOS, arch: universal}, Version{4, 0, 0, "rc1"}, "", ErrUnsupportedArch},

		{Platform{os: macOS, arch: amd64}, Version{4, 0, 0, "stable"}, "macos.universal", nil},
		{Platform{os: macOS, arch: arm64}, Version{4, 0, 0, "stable"}, "macos.universal", nil},
		{Platform{os: macOS, arch: i386}, Version{4, 0, 0, "stable"}, "", ErrUnsupportedArch},
		{Platform{os: macOS, arch: universal}, Version{4, 0, 0, "stable"}, "", ErrUnsupportedArch},

		// Valid inputs - Windows

		// v3.*
		{Platform{os: windows, arch: i386}, Version{major: 3}, "win32", nil},
		{Platform{os: windows, arch: amd64}, Version{major: 3}, "win64", nil},
		{Platform{os: windows, arch: arm64}, Version{major: 3}, "", ErrUnsupportedArch},
		{Platform{os: windows, arch: universal}, Version{major: 3}, "", ErrUnsupportedArch},

		// v4.0+
		{Platform{os: windows, arch: i386}, Version{major: 4}, "win32", nil},
		{Platform{os: windows, arch: amd64}, Version{major: 4}, "win64", nil},
		{Platform{os: windows, arch: arm64}, Version{major: 4}, "", ErrUnsupportedArch},
		{Platform{os: windows, arch: universal}, Version{major: 4}, "", ErrUnsupportedArch},
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
