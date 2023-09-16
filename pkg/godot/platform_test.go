package godot

import (
	"errors"
	"fmt"
	"testing"

	"github.com/coffeebeats/gdenv/internal/version"
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
		platform Platform
		version  string
		want     string
		err      error
	}{
		// Invalid inputs
		{platform: Platform{}, err: ErrMissingOS},
		{platform: Platform{os: linux}, err: ErrMissingArch},
		{platform: Platform{os: linux, arch: amd64}, err: version.ErrUnsupported},
		{platform: Platform{os: 100, arch: amd64}, err: ErrUnrecognizedOS},
		{platform: Platform{os: linux, arch: 100}, version: "3.0", err: ErrUnrecognizedArch},

		{platform: Platform{os: linux, arch: amd64}, version: "2.0", err: version.ErrUnsupported},
		{platform: Platform{os: macOS, arch: amd64}, version: "2.0", err: version.ErrUnsupported},
		{platform: Platform{os: windows, arch: amd64}, version: "2.0", err: version.ErrUnsupported},

		// Valid inputs - linux

		// v3.*
		{platform: Platform{os: linux, arch: i386}, version: "3.0", want: "x11.32"},
		{platform: Platform{os: linux, arch: amd64}, version: "3.0", want: "x11.64"},
		{platform: Platform{os: linux, arch: i386}, version: "3.0-stable_mono", want: "x11_32"},
		{platform: Platform{os: linux, arch: amd64}, version: "3.0-stable_mono", want: "x11_64"},
		{platform: Platform{os: linux, arch: arm64}, version: "3.0", err: ErrUnsupportedArch},
		{platform: Platform{os: linux, arch: universal}, version: "3.0", err: ErrUnsupportedArch},

		// v4.0-dev.* - v4.0-alpha14
		{platform: Platform{os: linux, arch: i386}, version: "4.0-dev.20220118", want: "linux.32"},
		{platform: Platform{os: linux, arch: amd64}, version: "4.0-dev.20220118", want: "linux.64"},
		{platform: Platform{os: linux, arch: arm64}, version: "4.0-dev.20220118", err: ErrUnsupportedArch},
		{platform: Platform{os: linux, arch: universal}, version: "4.0-dev.20220118", err: ErrUnsupportedArch},

		{platform: Platform{os: linux, arch: i386}, version: "4.0-alpha14", want: "linux.32"},
		{platform: Platform{os: linux, arch: amd64}, version: "4.0-alpha14", want: "linux.64"},
		{platform: Platform{os: linux, arch: arm64}, version: "4.0-alpha14", err: ErrUnsupportedArch},
		{platform: Platform{os: linux, arch: universal}, version: "4.0-alpha14", err: ErrUnsupportedArch},

		// v4.0-alpha15+
		{platform: Platform{os: linux, arch: i386}, version: "4.0-alpha15", want: "linux.x86_32"},
		{platform: Platform{os: linux, arch: amd64}, version: "4.0-alpha15", want: "linux.x86_64"},
		{platform: Platform{os: linux, arch: arm64}, version: "4.0-alpha15", err: ErrUnsupportedArch},
		{platform: Platform{os: linux, arch: universal}, version: "4.0-alpha15", err: ErrUnsupportedArch},

		{platform: Platform{os: linux, arch: i386}, version: "4.0", want: "linux.x86_32"},
		{platform: Platform{os: linux, arch: amd64}, version: "4.0", want: "linux.x86_64"},
		{platform: Platform{os: linux, arch: i386}, version: "4.0-stable_mono", want: "linux_x86_32"},
		{platform: Platform{os: linux, arch: amd64}, version: "4.0-stable_mono", want: "linux_x86_64"},
		{platform: Platform{os: linux, arch: arm64}, version: "4.0", err: ErrUnsupportedArch},
		{platform: Platform{os: linux, arch: universal}, version: "4.0", err: ErrUnsupportedArch},

		// Valid inputs - MacOS

		// v3.0 - v3.0.6
		{platform: Platform{os: macOS, arch: i386}, version: "3.0", want: "osx.fat"},
		{platform: Platform{os: macOS, arch: amd64}, version: "3.0", want: "osx.fat"},
		{platform: Platform{os: macOS, arch: arm64}, version: "3.0", err: ErrUnsupportedArch},
		{platform: Platform{os: macOS, arch: universal}, version: "3.0", err: ErrUnsupportedArch},

		{platform: Platform{os: macOS, arch: i386}, version: "3.0.6", want: "osx.fat"},
		{platform: Platform{os: macOS, arch: amd64}, version: "3.0.6", want: "osx.fat"},
		{platform: Platform{os: macOS, arch: arm64}, version: "3.0.6", err: ErrUnsupportedArch},
		{platform: Platform{os: macOS, arch: universal}, version: "3.0.6", err: ErrUnsupportedArch},

		// v3.1 - v3.2.4-beta2
		{platform: Platform{os: macOS, arch: amd64}, version: "3.1", want: "osx.64"},
		{platform: Platform{os: macOS, arch: i386}, version: "3.1", err: ErrUnsupportedArch},
		{platform: Platform{os: macOS, arch: arm64}, version: "3.1", err: ErrUnsupportedArch},
		{platform: Platform{os: macOS, arch: universal}, version: "3.1", err: ErrUnsupportedArch},

		{platform: Platform{os: macOS, arch: amd64}, version: "3.2.4-beta2", want: "osx.64"},
		{platform: Platform{os: macOS, arch: i386}, version: "3.2.4-beta2", err: ErrUnsupportedArch},
		{platform: Platform{os: macOS, arch: arm64}, version: "3.2.4-beta2", err: ErrUnsupportedArch},
		{platform: Platform{os: macOS, arch: universal}, version: "3.2.4-beta2", err: ErrUnsupportedArch},

		// v3.2.4-beta3 - v4.0-alpha12
		{platform: Platform{os: macOS, arch: amd64}, version: "3.2.4-beta3", want: "osx.universal"},
		{platform: Platform{os: macOS, arch: arm64}, version: "3.2.4-beta3", want: "osx.universal"},
		{platform: Platform{os: macOS, arch: i386}, version: "3.2.4-beta3", err: ErrUnsupportedArch},
		{platform: Platform{os: macOS, arch: universal}, version: "3.2.4-beta3", err: ErrUnsupportedArch},

		{platform: Platform{os: macOS, arch: amd64}, version: "3.2.4-rc1", want: "osx.universal"},
		{platform: Platform{os: macOS, arch: arm64}, version: "3.2.4-rc1", want: "osx.universal"},
		{platform: Platform{os: macOS, arch: i386}, version: "3.2.4-rc1", err: ErrUnsupportedArch},
		{platform: Platform{os: macOS, arch: universal}, version: "3.2.4-rc1", err: ErrUnsupportedArch},

		{platform: Platform{os: macOS, arch: amd64}, version: "3.2.4-stable", want: "osx.universal"},
		{platform: Platform{os: macOS, arch: arm64}, version: "3.2.4-stable", want: "osx.universal"},
		{platform: Platform{os: macOS, arch: i386}, version: "3.2.4-stable", err: ErrUnsupportedArch},
		{platform: Platform{os: macOS, arch: universal}, version: "3.2.4-stable", err: ErrUnsupportedArch},

		{platform: Platform{os: macOS, arch: amd64}, version: "4.0-dev.20210727", want: "osx.universal"},
		{platform: Platform{os: macOS, arch: arm64}, version: "4.0-dev.20210727", want: "osx.universal"},
		{platform: Platform{os: macOS, arch: i386}, version: "4.0-dev.20210727", err: ErrUnsupportedArch},
		{platform: Platform{os: macOS, arch: universal}, version: "4.0-dev.20210727", err: ErrUnsupportedArch},

		{platform: Platform{os: macOS, arch: amd64}, version: "4.0-alpha1", want: "osx.universal"},
		{platform: Platform{os: macOS, arch: arm64}, version: "4.0-alpha1", want: "osx.universal"},
		{platform: Platform{os: macOS, arch: i386}, version: "4.0-alpha1", err: ErrUnsupportedArch},
		{platform: Platform{os: macOS, arch: universal}, version: "4.0-alpha1", err: ErrUnsupportedArch},

		{platform: Platform{os: macOS, arch: amd64}, version: "4.0-alpha12", want: "osx.universal"},
		{platform: Platform{os: macOS, arch: arm64}, version: "4.0-alpha12", want: "osx.universal"},
		{platform: Platform{os: macOS, arch: i386}, version: "4.0-alpha12", err: ErrUnsupportedArch},
		{platform: Platform{os: macOS, arch: universal}, version: "4.0-alpha12", err: ErrUnsupportedArch},

		// v4.0-alpha13+
		{platform: Platform{os: macOS, arch: amd64}, version: "4.0-alpha13", want: "macos.universal"},
		{platform: Platform{os: macOS, arch: arm64}, version: "4.0-alpha13", want: "macos.universal"},
		{platform: Platform{os: macOS, arch: i386}, version: "4.0-alpha13", err: ErrUnsupportedArch},
		{platform: Platform{os: macOS, arch: universal}, version: "4.0-alpha13", err: ErrUnsupportedArch},

		{platform: Platform{os: macOS, arch: amd64}, version: "4.0-beta1", want: "macos.universal"},
		{platform: Platform{os: macOS, arch: arm64}, version: "4.0-beta1", want: "macos.universal"},
		{platform: Platform{os: macOS, arch: i386}, version: "4.0-beta1", err: ErrUnsupportedArch},
		{platform: Platform{os: macOS, arch: universal}, version: "4.0-beta1", err: ErrUnsupportedArch},

		{platform: Platform{os: macOS, arch: amd64}, version: "4.0-rc1", want: "macos.universal"},
		{platform: Platform{os: macOS, arch: arm64}, version: "4.0-rc1", want: "macos.universal"},
		{platform: Platform{os: macOS, arch: i386}, version: "4.0-rc1", err: ErrUnsupportedArch},
		{platform: Platform{os: macOS, arch: universal}, version: "4.0-rc1", err: ErrUnsupportedArch},

		{platform: Platform{os: macOS, arch: amd64}, version: "4.0", want: "macos.universal"},
		{platform: Platform{os: macOS, arch: arm64}, version: "4.0", want: "macos.universal"},
		{platform: Platform{os: macOS, arch: i386}, version: "4.0", err: ErrUnsupportedArch},
		{platform: Platform{os: macOS, arch: universal}, version: "4.0", err: ErrUnsupportedArch},

		// Valid inputs - Windows

		// v3.*
		{platform: Platform{os: windows, arch: i386}, version: "3.0", want: "win32"},
		{platform: Platform{os: windows, arch: amd64}, version: "3.0", want: "win64"},
		{platform: Platform{os: windows, arch: arm64}, version: "3.0", err: ErrUnsupportedArch},
		{platform: Platform{os: windows, arch: universal}, version: "3.0", err: ErrUnsupportedArch},

		// v4.0+
		{platform: Platform{os: windows, arch: i386}, version: "4.0", want: "win32"},
		{platform: Platform{os: windows, arch: amd64}, version: "4.0", want: "win64"},
		{platform: Platform{os: windows, arch: arm64}, version: "4.0", err: ErrUnsupportedArch},
		{platform: Platform{os: windows, arch: universal}, version: "4.0", err: ErrUnsupportedArch},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d-%v-%s", i, tc.platform, tc.version), func(t *testing.T) {
			var v version.Version
			if tc.version != "" {
				v = version.MustParse(tc.version)
			}

			got, err := FormatPlatform(tc.platform, v)

			if !errors.Is(err, tc.err) {
				t.Fatalf("err: got %#v, want %#v", err, tc.err)
			}
			if got != tc.want {
				t.Fatalf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}
