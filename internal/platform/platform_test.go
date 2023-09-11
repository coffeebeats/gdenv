package platform

import (
	"errors"
	"fmt"
	"testing"

	"github.com/coffeebeats/gdenv/pkg/godot"
)

/* --------------------------- Test: ParsePlatform -------------------------- */

func TestParsePlatform(t *testing.T) {
	tests := []struct {
		s    string
		want Platform
		err  error
	}{
		// Invalid inputs
		{s: "", err: godot.ErrMissingPlatform},
		{s: "abc", err: godot.ErrUnrecognizedPlatform},

		// Valid inputs (Go-defined)
		// Linux
		{s: "x11.32", want: Platform{I386, Linux}, err: nil},
		{s: "x11.64", want: Platform{Amd64, Linux}, err: nil},
		{s: "linux.x86_32", want: Platform{I386, Linux}, err: nil},
		{s: "linux.x86_64", want: Platform{Amd64, Linux}, err: nil},
		{s: "linux_x86_32", want: Platform{I386, Linux}, err: nil},
		{s: "linux_x86_64", want: Platform{Amd64, Linux}, err: nil},

		// MacOS
		{s: "osx.64", want: Platform{Amd64, MacOS}, err: nil},
		{s: "macos.universal", want: Platform{Universal, MacOS}, err: nil},
		{s: "osx.fat", want: Platform{Universal, MacOS}, err: nil},
		{s: "osx.universal", want: Platform{Universal, MacOS}, err: nil},

		// Windows
		{s: "win32", want: Platform{I386, Windows}, err: nil},
		{s: "win64", want: Platform{Amd64, Windows}, err: nil},

		// Valid inputs (user-supplied)
		{s: "WIN64", want: Platform{Amd64, Windows}, err: nil},
		{s: " MACOS.UNIVERSAL\n", want: Platform{Universal, MacOS}},
		{s: "\tlinux.x86_64 ", want: Platform{Amd64, Linux}, err: nil},
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
		v    godot.Version
		want string
		err  error
	}{
		// Invalid inputs
		{Platform{}, godot.Version{}, "", godot.ErrMissingOS},
		{Platform{OS: Linux}, godot.Version{}, "", godot.ErrMissingArch},
		{Platform{OS: Linux, Arch: Amd64}, godot.Version{}, "", godot.ErrUnsupportedVersion},
		{Platform{OS: 100, Arch: Amd64}, godot.Version{}, "", godot.ErrUnrecognizedOS},
		{Platform{OS: Linux, Arch: 100}, godot.Version{major: 3}, "", godot.ErrUnrecognizedArch},

		{Platform{OS: Linux, Arch: Amd64}, godot.Version{major: 2}, "", godot.ErrUnsupportedVersion},
		{Platform{OS: MacOS, Arch: Amd64}, godot.Version{major: 2}, "", godot.ErrUnsupportedVersion},
		{Platform{OS: Windows, Arch: Amd64}, godot.Version{major: 2}, "", godot.ErrUnsupportedVersion},

		// Valid inputs - linux

		// v3.*
		{Platform{OS: Linux, Arch: I386}, godot.Version{major: 3}, "x11.32", nil},
		{Platform{OS: Linux, Arch: Amd64}, godot.Version{major: 3}, "x11.64", nil},
		{Platform{OS: Linux, Arch: I386}, godot.Version{major: 3, label: LabelMono}, "x11_32", nil},
		{Platform{OS: Linux, Arch: Amd64}, godot.Version{major: 3, label: LabelMono}, "x11_64", nil},
		{Platform{OS: Linux, Arch: Arm64}, godot.Version{major: 3}, "", godot.ErrUnsupportedArch},
		{Platform{OS: Linux, Arch: Universal}, godot.Version{major: 3}, "", godot.ErrUnsupportedArch},

		// v4.0-dev.* - v4.0-alpha14
		{Platform{OS: Linux, Arch: I386}, godot.Version{major: 4, label: "dev.20220118"}, "linux.32", nil},
		{Platform{OS: Linux, Arch: Amd64}, godot.Version{major: 4, label: "dev.20220118"}, "linux.64", nil},
		{Platform{OS: Linux, Arch: Arm64}, godot.Version{major: 4, label: "dev.20220118"}, "", godot.ErrUnsupportedArch},
		{Platform{OS: Linux, Arch: Universal}, godot.Version{major: 4, label: "dev.20220118"}, "", godot.ErrUnsupportedArch},

		{Platform{OS: Linux, Arch: I386}, godot.Version{major: 4, label: "alpha14"}, "linux.32", nil},
		{Platform{OS: Linux, Arch: Amd64}, godot.Version{major: 4, label: "alpha14"}, "linux.64", nil},
		{Platform{OS: Linux, Arch: Arm64}, godot.Version{major: 4, label: "alpha14"}, "", godot.ErrUnsupportedArch},
		{Platform{OS: Linux, Arch: Universal}, godot.Version{major: 4, label: "alpha14"}, "", godot.ErrUnsupportedArch},

		// v4.0-alpha15+
		{Platform{OS: Linux, Arch: I386}, godot.Version{major: 4, label: "alpha15"}, "linux.x86_32", nil},
		{Platform{OS: Linux, Arch: Amd64}, godot.Version{major: 4, label: "alpha15"}, "linux.x86_64", nil},
		{Platform{OS: Linux, Arch: Arm64}, godot.Version{major: 4, label: "alpha15"}, "", godot.ErrUnsupportedArch},
		{Platform{OS: Linux, Arch: Universal}, godot.Version{major: 4, label: "alpha15"}, "", godot.ErrUnsupportedArch},

		{Platform{OS: Linux, Arch: I386}, godot.Version{major: 4}, "linux.x86_32", nil},
		{Platform{OS: Linux, Arch: Amd64}, godot.Version{major: 4}, "linux.x86_64", nil},
		{Platform{OS: Linux, Arch: I386}, godot.Version{major: 4, label: LabelMono}, "linux_x86_32", nil},
		{Platform{OS: Linux, Arch: Amd64}, godot.Version{major: 4, label: LabelMono}, "linux_x86_64", nil},
		{Platform{OS: Linux, Arch: Arm64}, godot.Version{major: 4}, "", godot.ErrUnsupportedArch},
		{Platform{OS: Linux, Arch: Universal}, godot.Version{major: 4}, "", godot.ErrUnsupportedArch},

		// Valid inputs - MacOS

		// v3.0 - v3.0.6
		{Platform{OS: MacOS, Arch: I386}, godot.Version{major: 3}, "osx.fat", nil},
		{Platform{OS: MacOS, Arch: Amd64}, godot.Version{major: 3}, "osx.fat", nil},
		{Platform{OS: MacOS, Arch: Arm64}, godot.Version{major: 3}, "", godot.ErrUnsupportedArch},
		{Platform{OS: MacOS, Arch: Universal}, godot.Version{major: 3}, "", godot.ErrUnsupportedArch},

		{Platform{OS: MacOS, Arch: I386}, godot.Version{major: 3, patch: 6}, "osx.fat", nil},
		{Platform{OS: MacOS, Arch: Amd64}, godot.Version{major: 3, patch: 6}, "osx.fat", nil},
		{Platform{OS: MacOS, Arch: Arm64}, godot.Version{major: 3, patch: 6}, "", godot.ErrUnsupportedArch},
		{Platform{OS: MacOS, Arch: Universal}, godot.Version{major: 3, patch: 6}, "", godot.ErrUnsupportedArch},

		// v3.1 - v3.2.4-beta2
		{Platform{OS: MacOS, Arch: Amd64}, godot.Version{major: 3, minor: 1}, "osx.64", nil},
		{Platform{OS: MacOS, Arch: I386}, godot.Version{major: 3, minor: 1}, "", godot.ErrUnsupportedArch},
		{Platform{OS: MacOS, Arch: Arm64}, godot.Version{major: 3, minor: 1}, "", godot.ErrUnsupportedArch},
		{Platform{OS: MacOS, Arch: Universal}, godot.Version{major: 3, minor: 1}, "", godot.ErrUnsupportedArch},

		{Platform{OS: MacOS, Arch: Amd64}, godot.Version{3, 2, 4, "beta2"}, "osx.64", nil},
		{Platform{OS: MacOS, Arch: I386}, godot.Version{3, 2, 4, "beta2"}, "", godot.ErrUnsupportedArch},
		{Platform{OS: MacOS, Arch: Arm64}, godot.Version{3, 2, 4, "beta2"}, "", godot.ErrUnsupportedArch},
		{Platform{OS: MacOS, Arch: Universal}, godot.Version{3, 2, 4, "beta2"}, "", godot.ErrUnsupportedArch},

		// v3.2.4-beta3 - v4.0-alpha12
		{Platform{OS: MacOS, Arch: Amd64}, godot.Version{3, 2, 4, "beta3"}, "osx.universal", nil},
		{Platform{OS: MacOS, Arch: Arm64}, godot.Version{3, 2, 4, "beta3"}, "osx.universal", nil},
		{Platform{OS: MacOS, Arch: I386}, godot.Version{3, 2, 4, "beta3"}, "", godot.ErrUnsupportedArch},
		{Platform{OS: MacOS, Arch: Universal}, godot.Version{3, 2, 4, "beta3"}, "", godot.ErrUnsupportedArch},

		{Platform{OS: MacOS, Arch: Amd64}, godot.Version{3, 2, 4, "rc1"}, "osx.universal", nil},
		{Platform{OS: MacOS, Arch: Arm64}, godot.Version{3, 2, 4, "rc1"}, "osx.universal", nil},
		{Platform{OS: MacOS, Arch: I386}, godot.Version{3, 2, 4, "rc1"}, "", godot.ErrUnsupportedArch},
		{Platform{OS: MacOS, Arch: Universal}, godot.Version{3, 2, 4, "rc1"}, "", godot.ErrUnsupportedArch},

		{Platform{OS: MacOS, Arch: Amd64}, godot.Version{3, 2, 4, "stable"}, "osx.universal", nil},
		{Platform{OS: MacOS, Arch: Arm64}, godot.Version{3, 2, 4, "stable"}, "osx.universal", nil},
		{Platform{OS: MacOS, Arch: I386}, godot.Version{3, 2, 4, "stable"}, "", godot.ErrUnsupportedArch},
		{Platform{OS: MacOS, Arch: Universal}, godot.Version{3, 2, 4, "stable"}, "", godot.ErrUnsupportedArch},

		{Platform{OS: MacOS, Arch: Amd64}, godot.Version{4, 0, 0, "dev.20210727"}, "osx.universal", nil},
		{Platform{OS: MacOS, Arch: Arm64}, godot.Version{4, 0, 0, "dev.20210727"}, "osx.universal", nil},
		{Platform{OS: MacOS, Arch: I386}, godot.Version{4, 0, 0, "dev.20210727"}, "", godot.ErrUnsupportedArch},
		{Platform{OS: MacOS, Arch: Universal}, godot.Version{4, 0, 0, "dev.20210727"}, "", godot.ErrUnsupportedArch},

		{Platform{OS: MacOS, Arch: Amd64}, godot.Version{4, 0, 0, "alpha1"}, "osx.universal", nil},
		{Platform{OS: MacOS, Arch: Arm64}, godot.Version{4, 0, 0, "alpha1"}, "osx.universal", nil},
		{Platform{OS: MacOS, Arch: I386}, godot.Version{4, 0, 0, "alpha1"}, "", godot.ErrUnsupportedArch},
		{Platform{OS: MacOS, Arch: Universal}, godot.Version{4, 0, 0, "alpha1"}, "", godot.ErrUnsupportedArch},

		{Platform{OS: MacOS, Arch: Amd64}, godot.Version{4, 0, 0, "alpha12"}, "osx.universal", nil},
		{Platform{OS: MacOS, Arch: Arm64}, godot.Version{4, 0, 0, "alpha12"}, "osx.universal", nil},
		{Platform{OS: MacOS, Arch: I386}, godot.Version{4, 0, 0, "alpha12"}, "", godot.ErrUnsupportedArch},
		{Platform{OS: MacOS, Arch: Universal}, godot.Version{4, 0, 0, "alpha12"}, "", godot.ErrUnsupportedArch},

		// v4.0-alpha13+
		{Platform{OS: MacOS, Arch: Amd64}, godot.Version{4, 0, 0, "alpha13"}, "macos.universal", nil},
		{Platform{OS: MacOS, Arch: Arm64}, godot.Version{4, 0, 0, "alpha13"}, "macos.universal", nil},
		{Platform{OS: MacOS, Arch: I386}, godot.Version{4, 0, 0, "alpha13"}, "", godot.ErrUnsupportedArch},
		{Platform{OS: MacOS, Arch: Universal}, godot.Version{4, 0, 0, "alpha13"}, "", godot.ErrUnsupportedArch},

		{Platform{OS: MacOS, Arch: Amd64}, godot.Version{4, 0, 0, "beta1"}, "macos.universal", nil},
		{Platform{OS: MacOS, Arch: Arm64}, godot.Version{4, 0, 0, "beta1"}, "macos.universal", nil},
		{Platform{OS: MacOS, Arch: I386}, godot.Version{4, 0, 0, "beta1"}, "", godot.ErrUnsupportedArch},
		{Platform{OS: MacOS, Arch: Universal}, godot.Version{4, 0, 0, "beta1"}, "", godot.ErrUnsupportedArch},

		{Platform{OS: MacOS, Arch: Amd64}, godot.Version{4, 0, 0, "rc1"}, "macos.universal", nil},
		{Platform{OS: MacOS, Arch: Arm64}, godot.Version{4, 0, 0, "rc1"}, "macos.universal", nil},
		{Platform{OS: MacOS, Arch: I386}, godot.Version{4, 0, 0, "rc1"}, "", godot.ErrUnsupportedArch},
		{Platform{OS: MacOS, Arch: Universal}, godot.Version{4, 0, 0, "rc1"}, "", godot.ErrUnsupportedArch},

		{Platform{OS: MacOS, Arch: Amd64}, godot.Version{4, 0, 0, "stable"}, "macos.universal", nil},
		{Platform{OS: MacOS, Arch: Arm64}, godot.Version{4, 0, 0, "stable"}, "macos.universal", nil},
		{Platform{OS: MacOS, Arch: I386}, godot.Version{4, 0, 0, "stable"}, "", godot.ErrUnsupportedArch},
		{Platform{OS: MacOS, Arch: Universal}, godot.Version{4, 0, 0, "stable"}, "", godot.ErrUnsupportedArch},

		// Valid inputs - Windows

		// v3.*
		{Platform{OS: Windows, Arch: I386}, Version{major: 3}, "win32", nil},
		{Platform{OS: Windows, Arch: Amd64}, Version{major: 3}, "win64", nil},
		{Platform{OS: Windows, Arch: Arm64}, Version{major: 3}, "", godot.ErrUnsupportedArch},
		{Platform{OS: Windows, Arch: Universal}, Version{major: 3}, "", godot.ErrUnsupportedArch},

		// v4.0+
		{Platform{OS: Windows, Arch: I386}, Version{major: 4}, "win32", nil},
		{Platform{OS: Windows, Arch: Amd64}, Version{major: 4}, "win64", nil},
		{Platform{OS: Windows, Arch: Arm64}, Version{major: 4}, "", godot.ErrUnsupportedArch},
		{Platform{OS: Windows, Arch: Universal}, Version{major: 4}, "", godot.ErrUnsupportedArch},
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
