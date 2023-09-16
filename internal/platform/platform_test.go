package platform

import (
	"errors"
	"fmt"
	"testing"

	"github.com/coffeebeats/gdenv/internal/version"
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
			got, err := Parse(tc.s)

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
		{platform: Platform{OS: Linux}, err: ErrMissingArch},
		{platform: Platform{OS: Linux, Arch: Amd64}, err: version.ErrUnsupported},
		{platform: Platform{OS: 100, Arch: Amd64}, err: ErrUnrecognizedOS},
		{platform: Platform{OS: Linux, Arch: 100}, version: "3.0", err: ErrUnrecognizedArch},

		{platform: Platform{OS: Linux, Arch: Amd64}, version: "2.0", err: version.ErrUnsupported},
		{platform: Platform{OS: MacOS, Arch: Amd64}, version: "2.0", err: version.ErrUnsupported},
		{platform: Platform{OS: Windows, Arch: Amd64}, version: "2.0", err: version.ErrUnsupported},

		// Valid inputs - linux

		// v3.*
		{platform: Platform{OS: Linux, Arch: I386}, version: "3.0", want: "x11.32"},
		{platform: Platform{OS: Linux, Arch: Amd64}, version: "3.0", want: "x11.64"},
		{platform: Platform{OS: Linux, Arch: I386}, version: "3.0-stable_mono", want: "x11_32"},
		{platform: Platform{OS: Linux, Arch: Amd64}, version: "3.0-stable_mono", want: "x11_64"},
		{platform: Platform{OS: Linux, Arch: Arm64}, version: "3.0", err: ErrUnrecognizedArch},
		{platform: Platform{OS: Linux, Arch: Universal}, version: "3.0", err: ErrUnrecognizedArch},

		// v4.0-dev.* - v4.0-alpha14
		{platform: Platform{OS: Linux, Arch: I386}, version: "4.0-dev.20220118", want: "linux.32"},
		{platform: Platform{OS: Linux, Arch: Amd64}, version: "4.0-dev.20220118", want: "linux.64"},
		{platform: Platform{OS: Linux, Arch: Arm64}, version: "4.0-dev.20220118", err: ErrUnrecognizedArch},
		{platform: Platform{OS: Linux, Arch: Universal}, version: "4.0-dev.20220118", err: ErrUnrecognizedArch},

		{platform: Platform{OS: Linux, Arch: I386}, version: "4.0-alpha14", want: "linux.32"},
		{platform: Platform{OS: Linux, Arch: Amd64}, version: "4.0-alpha14", want: "linux.64"},
		{platform: Platform{OS: Linux, Arch: Arm64}, version: "4.0-alpha14", err: ErrUnrecognizedArch},
		{platform: Platform{OS: Linux, Arch: Universal}, version: "4.0-alpha14", err: ErrUnrecognizedArch},

		// v4.0-alpha15+
		{platform: Platform{OS: Linux, Arch: I386}, version: "4.0-alpha15", want: "linux.x86_32"},
		{platform: Platform{OS: Linux, Arch: Amd64}, version: "4.0-alpha15", want: "linux.x86_64"},
		{platform: Platform{OS: Linux, Arch: Arm64}, version: "4.0-alpha15", err: ErrUnrecognizedArch},
		{platform: Platform{OS: Linux, Arch: Universal}, version: "4.0-alpha15", err: ErrUnrecognizedArch},

		{platform: Platform{OS: Linux, Arch: I386}, version: "4.0", want: "linux.x86_32"},
		{platform: Platform{OS: Linux, Arch: Amd64}, version: "4.0", want: "linux.x86_64"},
		{platform: Platform{OS: Linux, Arch: I386}, version: "4.0-stable_mono", want: "linux_x86_32"},
		{platform: Platform{OS: Linux, Arch: Amd64}, version: "4.0-stable_mono", want: "linux_x86_64"},
		{platform: Platform{OS: Linux, Arch: Arm64}, version: "4.0", err: ErrUnrecognizedArch},
		{platform: Platform{OS: Linux, Arch: Universal}, version: "4.0", err: ErrUnrecognizedArch},

		// Valid inputs - MacOS

		// v3.0 - v3.0.6
		{platform: Platform{OS: MacOS, Arch: I386}, version: "3.0", want: "osx.fat"},
		{platform: Platform{OS: MacOS, Arch: Amd64}, version: "3.0", want: "osx.fat"},
		{platform: Platform{OS: MacOS, Arch: Arm64}, version: "3.0", err: ErrUnrecognizedArch},
		{platform: Platform{OS: MacOS, Arch: Universal}, version: "3.0", err: ErrUnrecognizedArch},

		{platform: Platform{OS: MacOS, Arch: I386}, version: "3.0.6", want: "osx.fat"},
		{platform: Platform{OS: MacOS, Arch: Amd64}, version: "3.0.6", want: "osx.fat"},
		{platform: Platform{OS: MacOS, Arch: Arm64}, version: "3.0.6", err: ErrUnrecognizedArch},
		{platform: Platform{OS: MacOS, Arch: Universal}, version: "3.0.6", err: ErrUnrecognizedArch},

		// v3.1 - v3.2.4-beta2
		{platform: Platform{OS: MacOS, Arch: Amd64}, version: "3.1", want: "osx.64"},
		{platform: Platform{OS: MacOS, Arch: I386}, version: "3.1", err: ErrUnrecognizedArch},
		{platform: Platform{OS: MacOS, Arch: Arm64}, version: "3.1", err: ErrUnrecognizedArch},
		{platform: Platform{OS: MacOS, Arch: Universal}, version: "3.1", err: ErrUnrecognizedArch},

		{platform: Platform{OS: MacOS, Arch: Amd64}, version: "3.2.4-beta2", want: "osx.64"},
		{platform: Platform{OS: MacOS, Arch: I386}, version: "3.2.4-beta2", err: ErrUnrecognizedArch},
		{platform: Platform{OS: MacOS, Arch: Arm64}, version: "3.2.4-beta2", err: ErrUnrecognizedArch},
		{platform: Platform{OS: MacOS, Arch: Universal}, version: "3.2.4-beta2", err: ErrUnrecognizedArch},

		// v3.2.4-beta3 - v4.0-alpha12
		{platform: Platform{OS: MacOS, Arch: Amd64}, version: "3.2.4-beta3", want: "osx.universal"},
		{platform: Platform{OS: MacOS, Arch: Arm64}, version: "3.2.4-beta3", want: "osx.universal"},
		{platform: Platform{OS: MacOS, Arch: I386}, version: "3.2.4-beta3", err: ErrUnrecognizedArch},
		{platform: Platform{OS: MacOS, Arch: Universal}, version: "3.2.4-beta3", err: ErrUnrecognizedArch},

		{platform: Platform{OS: MacOS, Arch: Amd64}, version: "3.2.4-rc1", want: "osx.universal"},
		{platform: Platform{OS: MacOS, Arch: Arm64}, version: "3.2.4-rc1", want: "osx.universal"},
		{platform: Platform{OS: MacOS, Arch: I386}, version: "3.2.4-rc1", err: ErrUnrecognizedArch},
		{platform: Platform{OS: MacOS, Arch: Universal}, version: "3.2.4-rc1", err: ErrUnrecognizedArch},

		{platform: Platform{OS: MacOS, Arch: Amd64}, version: "3.2.4-stable", want: "osx.universal"},
		{platform: Platform{OS: MacOS, Arch: Arm64}, version: "3.2.4-stable", want: "osx.universal"},
		{platform: Platform{OS: MacOS, Arch: I386}, version: "3.2.4-stable", err: ErrUnrecognizedArch},
		{platform: Platform{OS: MacOS, Arch: Universal}, version: "3.2.4-stable", err: ErrUnrecognizedArch},

		{platform: Platform{OS: MacOS, Arch: Amd64}, version: "4.0-dev.20210727", want: "osx.universal"},
		{platform: Platform{OS: MacOS, Arch: Arm64}, version: "4.0-dev.20210727", want: "osx.universal"},
		{platform: Platform{OS: MacOS, Arch: I386}, version: "4.0-dev.20210727", err: ErrUnrecognizedArch},
		{platform: Platform{OS: MacOS, Arch: Universal}, version: "4.0-dev.20210727", err: ErrUnrecognizedArch},

		{platform: Platform{OS: MacOS, Arch: Amd64}, version: "4.0-alpha1", want: "osx.universal"},
		{platform: Platform{OS: MacOS, Arch: Arm64}, version: "4.0-alpha1", want: "osx.universal"},
		{platform: Platform{OS: MacOS, Arch: I386}, version: "4.0-alpha1", err: ErrUnrecognizedArch},
		{platform: Platform{OS: MacOS, Arch: Universal}, version: "4.0-alpha1", err: ErrUnrecognizedArch},

		{platform: Platform{OS: MacOS, Arch: Amd64}, version: "4.0-alpha12", want: "osx.universal"},
		{platform: Platform{OS: MacOS, Arch: Arm64}, version: "4.0-alpha12", want: "osx.universal"},
		{platform: Platform{OS: MacOS, Arch: I386}, version: "4.0-alpha12", err: ErrUnrecognizedArch},
		{platform: Platform{OS: MacOS, Arch: Universal}, version: "4.0-alpha12", err: ErrUnrecognizedArch},

		// v4.0-alpha13+
		{platform: Platform{OS: MacOS, Arch: Amd64}, version: "4.0-alpha13", want: "macos.universal"},
		{platform: Platform{OS: MacOS, Arch: Arm64}, version: "4.0-alpha13", want: "macos.universal"},
		{platform: Platform{OS: MacOS, Arch: I386}, version: "4.0-alpha13", err: ErrUnrecognizedArch},
		{platform: Platform{OS: MacOS, Arch: Universal}, version: "4.0-alpha13", err: ErrUnrecognizedArch},

		{platform: Platform{OS: MacOS, Arch: Amd64}, version: "4.0-beta1", want: "macos.universal"},
		{platform: Platform{OS: MacOS, Arch: Arm64}, version: "4.0-beta1", want: "macos.universal"},
		{platform: Platform{OS: MacOS, Arch: I386}, version: "4.0-beta1", err: ErrUnrecognizedArch},
		{platform: Platform{OS: MacOS, Arch: Universal}, version: "4.0-beta1", err: ErrUnrecognizedArch},

		{platform: Platform{OS: MacOS, Arch: Amd64}, version: "4.0-rc1", want: "macos.universal"},
		{platform: Platform{OS: MacOS, Arch: Arm64}, version: "4.0-rc1", want: "macos.universal"},
		{platform: Platform{OS: MacOS, Arch: I386}, version: "4.0-rc1", err: ErrUnrecognizedArch},
		{platform: Platform{OS: MacOS, Arch: Universal}, version: "4.0-rc1", err: ErrUnrecognizedArch},

		{platform: Platform{OS: MacOS, Arch: Amd64}, version: "4.0", want: "macos.universal"},
		{platform: Platform{OS: MacOS, Arch: Arm64}, version: "4.0", want: "macos.universal"},
		{platform: Platform{OS: MacOS, Arch: I386}, version: "4.0", err: ErrUnrecognizedArch},
		{platform: Platform{OS: MacOS, Arch: Universal}, version: "4.0", err: ErrUnrecognizedArch},

		// Valid inputs - Windows

		// v3.*
		{platform: Platform{OS: Windows, Arch: I386}, version: "3.0", want: "win32"},
		{platform: Platform{OS: Windows, Arch: Amd64}, version: "3.0", want: "win64"},
		{platform: Platform{OS: Windows, Arch: Arm64}, version: "3.0", err: ErrUnrecognizedArch},
		{platform: Platform{OS: Windows, Arch: Universal}, version: "3.0", err: ErrUnrecognizedArch},

		// v4.0+
		{platform: Platform{OS: Windows, Arch: I386}, version: "4.0", want: "win32"},
		{platform: Platform{OS: Windows, Arch: Amd64}, version: "4.0", want: "win64"},
		{platform: Platform{OS: Windows, Arch: Arm64}, version: "4.0", err: ErrUnrecognizedArch},
		{platform: Platform{OS: Windows, Arch: Universal}, version: "4.0", err: ErrUnrecognizedArch},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d-%v-%s", i, tc.platform, tc.version), func(t *testing.T) {
			var v version.Version
			if tc.version != "" {
				v = version.MustParse(tc.version)
			}

			got, err := Format(tc.platform, v)

			if !errors.Is(err, tc.err) {
				t.Fatalf("err: got %#v, want %#v", err, tc.err)
			}
			if got != tc.want {
				t.Fatalf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}
