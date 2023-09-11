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
		{s: "", err: ErrMissingPlatform},
		{s: "abc", err: ErrUnrecognizedPlatform},

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
	var (
		v2 = godot.NewVersion(2)
		v3 = godot.NewVersion(3)
		v4 = godot.NewVersion(4)
	)

	tests := []struct {
		p    Platform
		v    godot.Version
		want string
		err  error
	}{
		// Invalid inputs
		{Platform{}, godot.Version{}, "", ErrMissingOS},
		{Platform{OS: Linux}, godot.Version{}, "", ErrMissingArch},
		{Platform{OS: Linux, Arch: Amd64}, godot.Version{}, "", godot.ErrUnsupportedVersion},
		{Platform{OS: 100, Arch: Amd64}, godot.Version{}, "", ErrUnrecognizedOS},
		{Platform{OS: Linux, Arch: 100}, v3, "", ErrUnrecognizedArch},

		{Platform{OS: Linux, Arch: Amd64}, v2, "", godot.ErrUnsupportedVersion},
		{Platform{OS: MacOS, Arch: Amd64}, v2, "", godot.ErrUnsupportedVersion},
		{Platform{OS: Windows, Arch: Amd64}, v2, "", godot.ErrUnsupportedVersion},

		// Valid inputs - linux

		// v3.*
		{Platform{OS: Linux, Arch: I386}, v3, "x11.32", nil},
		{Platform{OS: Linux, Arch: Amd64}, v3, "x11.64", nil},
		{Platform{OS: Linux, Arch: I386}, v3.WithLabel(godot.LabelMono), "x11_32", nil},
		{Platform{OS: Linux, Arch: Amd64}, v3.WithLabel(godot.LabelMono), "x11_64", nil},
		{Platform{OS: Linux, Arch: Arm64}, v3, "", ErrUnsupportedArch},
		{Platform{OS: Linux, Arch: Universal}, v3, "", ErrUnsupportedArch},

		// v4.0-dev.* - v4.0-alpha14
		{Platform{OS: Linux, Arch: I386}, v4.WithLabel("dev.20220118"), "linux.32", nil},
		{Platform{OS: Linux, Arch: Amd64}, v4.WithLabel("dev.20220118"), "linux.64", nil},
		{Platform{OS: Linux, Arch: Arm64}, v4.WithLabel("dev.20220118"), "", ErrUnsupportedArch},
		{Platform{OS: Linux, Arch: Universal}, v4.WithLabel("dev.20220118"), "", ErrUnsupportedArch},

		{Platform{OS: Linux, Arch: I386}, v4.WithLabel("alpha14"), "linux.32", nil},
		{Platform{OS: Linux, Arch: Amd64}, v4.WithLabel("alpha14"), "linux.64", nil},
		{Platform{OS: Linux, Arch: Arm64}, v4.WithLabel("alpha14"), "", ErrUnsupportedArch},
		{Platform{OS: Linux, Arch: Universal}, v4.WithLabel("alpha14"), "", ErrUnsupportedArch},

		// v4.0-alpha15+
		{Platform{OS: Linux, Arch: I386}, v4.WithLabel("alpha15"), "linux.x86_32", nil},
		{Platform{OS: Linux, Arch: Amd64}, v4.WithLabel("alpha15"), "linux.x86_64", nil},
		{Platform{OS: Linux, Arch: Arm64}, v4.WithLabel("alpha15"), "", ErrUnsupportedArch},
		{Platform{OS: Linux, Arch: Universal}, v4.WithLabel("alpha15"), "", ErrUnsupportedArch},

		{Platform{OS: Linux, Arch: I386}, v4, "linux.x86_32", nil},
		{Platform{OS: Linux, Arch: Amd64}, v4, "linux.x86_64", nil},
		{Platform{OS: Linux, Arch: I386}, v4.WithLabel(godot.LabelMono), "linux_x86_32", nil},
		{Platform{OS: Linux, Arch: Amd64}, v4.WithLabel(godot.LabelMono), "linux_x86_64", nil},
		{Platform{OS: Linux, Arch: Arm64}, v4, "", ErrUnsupportedArch},
		{Platform{OS: Linux, Arch: Universal}, v4, "", ErrUnsupportedArch},

		// Valid inputs - MacOS

		// v3.0 - v3.0.6
		{Platform{OS: MacOS, Arch: I386}, v3, "osx.fat", nil},
		{Platform{OS: MacOS, Arch: Amd64}, v3, "osx.fat", nil},
		{Platform{OS: MacOS, Arch: Arm64}, v3, "", ErrUnsupportedArch},
		{Platform{OS: MacOS, Arch: Universal}, v3, "", ErrUnsupportedArch},

		{Platform{OS: MacOS, Arch: I386}, godot.NewVersion(3, 0, 6), "osx.fat", nil},
		{Platform{OS: MacOS, Arch: Amd64}, godot.NewVersion(3, 0, 6), "osx.fat", nil},
		{Platform{OS: MacOS, Arch: Arm64}, godot.NewVersion(3, 0, 6), "", ErrUnsupportedArch},
		{Platform{OS: MacOS, Arch: Universal}, godot.NewVersion(3, 0, 6), "", ErrUnsupportedArch},

		// v3.1 - v3.2.4-beta2
		{Platform{OS: MacOS, Arch: Amd64}, godot.NewVersion(3, 1), "osx.64", nil},
		{Platform{OS: MacOS, Arch: I386}, godot.NewVersion(3, 1), "", ErrUnsupportedArch},
		{Platform{OS: MacOS, Arch: Arm64}, godot.NewVersion(3, 1), "", ErrUnsupportedArch},
		{Platform{OS: MacOS, Arch: Universal}, godot.NewVersion(3, 1), "", ErrUnsupportedArch},

		{Platform{OS: MacOS, Arch: Amd64}, godot.NewVersionWithLabel(3, 2, 4, "beta2"), "osx.64", nil},
		{Platform{OS: MacOS, Arch: I386}, godot.NewVersionWithLabel(3, 2, 4, "beta2"), "", ErrUnsupportedArch},
		{Platform{OS: MacOS, Arch: Arm64}, godot.NewVersionWithLabel(3, 2, 4, "beta2"), "", ErrUnsupportedArch},
		{Platform{OS: MacOS, Arch: Universal}, godot.NewVersionWithLabel(3, 2, 4, "beta2"), "", ErrUnsupportedArch},

		// v3.2.4-beta3 - v4.0-alpha12
		{Platform{OS: MacOS, Arch: Amd64}, godot.NewVersionWithLabel(3, 2, 4, "beta3"), "osx.universal", nil},
		{Platform{OS: MacOS, Arch: Arm64}, godot.NewVersionWithLabel(3, 2, 4, "beta3"), "osx.universal", nil},
		{Platform{OS: MacOS, Arch: I386}, godot.NewVersionWithLabel(3, 2, 4, "beta3"), "", ErrUnsupportedArch},
		{Platform{OS: MacOS, Arch: Universal}, godot.NewVersionWithLabel(3, 2, 4, "beta3"), "", ErrUnsupportedArch},

		{Platform{OS: MacOS, Arch: Amd64}, godot.NewVersionWithLabel(3, 2, 4, "rc1"), "osx.universal", nil},
		{Platform{OS: MacOS, Arch: Arm64}, godot.NewVersionWithLabel(3, 2, 4, "rc1"), "osx.universal", nil},
		{Platform{OS: MacOS, Arch: I386}, godot.NewVersionWithLabel(3, 2, 4, "rc1"), "", ErrUnsupportedArch},
		{Platform{OS: MacOS, Arch: Universal}, godot.NewVersionWithLabel(3, 2, 4, "rc1"), "", ErrUnsupportedArch},

		{Platform{OS: MacOS, Arch: Amd64}, godot.NewVersion(3, 2, 4), "osx.universal", nil},
		{Platform{OS: MacOS, Arch: Arm64}, godot.NewVersion(3, 2, 4), "osx.universal", nil},
		{Platform{OS: MacOS, Arch: I386}, godot.NewVersion(3, 2, 4), "", ErrUnsupportedArch},
		{Platform{OS: MacOS, Arch: Universal}, godot.NewVersion(3, 2, 4), "", ErrUnsupportedArch},

		{Platform{OS: MacOS, Arch: Amd64}, v4.WithLabel("dev.20210727"), "osx.universal", nil},
		{Platform{OS: MacOS, Arch: Arm64}, v4.WithLabel("dev.20210727"), "osx.universal", nil},
		{Platform{OS: MacOS, Arch: I386}, v4.WithLabel("dev.20210727"), "", ErrUnsupportedArch},
		{Platform{OS: MacOS, Arch: Universal}, v4.WithLabel("dev.20210727"), "", ErrUnsupportedArch},

		{Platform{OS: MacOS, Arch: Amd64}, v4.WithLabel("alpha1"), "osx.universal", nil},
		{Platform{OS: MacOS, Arch: Arm64}, v4.WithLabel("alpha1"), "osx.universal", nil},
		{Platform{OS: MacOS, Arch: I386}, v4.WithLabel("alpha1"), "", ErrUnsupportedArch},
		{Platform{OS: MacOS, Arch: Universal}, v4.WithLabel("alpha1"), "", ErrUnsupportedArch},

		{Platform{OS: MacOS, Arch: Amd64}, v4.WithLabel("alpha12"), "osx.universal", nil},
		{Platform{OS: MacOS, Arch: Arm64}, v4.WithLabel("alpha12"), "osx.universal", nil},
		{Platform{OS: MacOS, Arch: I386}, v4.WithLabel("alpha12"), "", ErrUnsupportedArch},
		{Platform{OS: MacOS, Arch: Universal}, v4.WithLabel("alpha12"), "", ErrUnsupportedArch},

		// v4.0-alpha13+
		{Platform{OS: MacOS, Arch: Amd64}, v4.WithLabel("alpha13"), "macos.universal", nil},
		{Platform{OS: MacOS, Arch: Arm64}, v4.WithLabel("alpha13"), "macos.universal", nil},
		{Platform{OS: MacOS, Arch: I386}, v4.WithLabel("alpha13"), "", ErrUnsupportedArch},
		{Platform{OS: MacOS, Arch: Universal}, v4.WithLabel("alpha13"), "", ErrUnsupportedArch},

		{Platform{OS: MacOS, Arch: Amd64}, v4.WithLabel("beta1"), "macos.universal", nil},
		{Platform{OS: MacOS, Arch: Arm64}, v4.WithLabel("beta1"), "macos.universal", nil},
		{Platform{OS: MacOS, Arch: I386}, v4.WithLabel("beta1"), "", ErrUnsupportedArch},
		{Platform{OS: MacOS, Arch: Universal}, v4.WithLabel("beta1"), "", ErrUnsupportedArch},

		{Platform{OS: MacOS, Arch: Amd64}, v4.WithLabel("rc1"), "macos.universal", nil},
		{Platform{OS: MacOS, Arch: Arm64}, v4.WithLabel("rc1"), "macos.universal", nil},
		{Platform{OS: MacOS, Arch: I386}, v4.WithLabel("rc1"), "", ErrUnsupportedArch},
		{Platform{OS: MacOS, Arch: Universal}, v4.WithLabel("rc1"), "", ErrUnsupportedArch},

		{Platform{OS: MacOS, Arch: Amd64}, v4, "macos.universal", nil},
		{Platform{OS: MacOS, Arch: Arm64}, v4, "macos.universal", nil},
		{Platform{OS: MacOS, Arch: I386}, v4, "", ErrUnsupportedArch},
		{Platform{OS: MacOS, Arch: Universal}, v4, "", ErrUnsupportedArch},

		// Valid inputs - Windows

		// v3.*
		{Platform{OS: Windows, Arch: I386}, v3, "win32", nil},
		{Platform{OS: Windows, Arch: Amd64}, v3, "win64", nil},
		{Platform{OS: Windows, Arch: Arm64}, v3, "", ErrUnsupportedArch},
		{Platform{OS: Windows, Arch: Universal}, v3, "", ErrUnsupportedArch},

		// v4.0+
		{Platform{OS: Windows, Arch: I386}, v4, "win32", nil},
		{Platform{OS: Windows, Arch: Amd64}, v4, "win64", nil},
		{Platform{OS: Windows, Arch: Arm64}, v4, "", ErrUnsupportedArch},
		{Platform{OS: Windows, Arch: Universal}, v4, "", ErrUnsupportedArch},
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
