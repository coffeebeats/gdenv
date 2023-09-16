package godot

import (
	"errors"
	"fmt"
	"testing"

	"github.com/coffeebeats/gdenv/internal/godot/platform"
	"github.com/coffeebeats/gdenv/internal/godot/version"
)

/* -------------------------- Test: Executable.Name ------------------------- */

func TestExecutableName(t *testing.T) {
	tests := []struct {
		platform platform.Platform
		version  string
		want     string
		err      error
	}{
		// Invalid inputs
		{platform: platform.Platform{}, err: platform.ErrMissingOS},
		{platform: platform.Platform{OS: platform.Linux}, err: platform.ErrMissingArch},
		{platform: linux64(), err: version.ErrUnsupported},

		{platform: platform.Platform{OS: platform.Linux, Arch: platform.Universal}, version: "4.0", err: platform.ErrUnrecognizedArch},

		{platform: linux64(), version: "2.0", err: version.ErrUnsupported},
		{platform: macOSX86_64(), version: "2.0", err: version.ErrUnsupported},
		{platform: windows64(), version: "2.0", err: version.ErrUnsupported},

		// Valid inputs - Linux

		// v3.6-beta1
		{platform: linux32(), version: "3.6-beta1", want: "Godot_v3.6-beta1_x11.32"},
		{platform: linux64(), version: "3.6-beta1", want: "Godot_v3.6-beta1_x11.64"},

		// v4.0
		{platform: linux32(), version: "4.0", want: "Godot_v4.0-stable_linux.x86_32"},
		{platform: linux64(), version: "4.0", want: "Godot_v4.0-stable_linux.x86_64"},

		// v4.0-stable_mono
		{platform: linux32(), version: "4.0-stable_mono", want: "Godot_v4.0-stable_mono_linux_x86_32"},
		{platform: linux64(), version: "4.0-stable_mono", want: "Godot_v4.0-stable_mono_linux_x86_64"},

		// v5.0-rc4
		{platform: linux32(), version: "5.0-rc4", want: "Godot_v5.0-rc4_linux.x86_32"},
		{platform: linux64(), version: "5.0-rc4", want: "Godot_v5.0-rc4_linux.x86_64"},

		// Valid inputs - MacOS

		// v3.6-beta1
		{platform: macOSX86_64(), version: "3.6-beta1", want: "Godot_v3.6-beta1_osx.universal"},
		{platform: macOSArm64(), version: "3.6-beta1", want: "Godot_v3.6-beta1_osx.universal"},

		// v4.0
		{platform: macOSX86_64(), version: "4.0", want: "Godot_v4.0-stable_macos.universal"},
		{platform: macOSArm64(), version: "4.0", want: "Godot_v4.0-stable_macos.universal"},

		// v4.0-stable_mono
		{platform: macOSX86_64(), version: "4.0-stable_mono", want: "Godot_v4.0-stable_mono_macos.universal"},
		{platform: macOSArm64(), version: "4.0-stable_mono", want: "Godot_v4.0-stable_mono_macos.universal"},

		// v5.0-rc4
		{platform: macOSX86_64(), version: "5.0-rc4", want: "Godot_v5.0-rc4_macos.universal"},
		{platform: macOSArm64(), version: "5.0-rc4", want: "Godot_v5.0-rc4_macos.universal"},

		// Valid inputs - Windows

		// v3.6-beta1
		{platform: platform.Platform{OS: platform.Windows, Arch: platform.I386}, version: "3.6-beta1", want: "Godot_v3.6-beta1_win32.exe"},
		{platform: platform.Platform{OS: platform.Windows, Arch: platform.Amd64}, version: "3.6-beta1", want: "Godot_v3.6-beta1_win64.exe"},

		// v4.0
		{platform: platform.Platform{OS: platform.Windows, Arch: platform.I386}, version: "4.0", want: "Godot_v4.0-stable_win32.exe"},
		{platform: platform.Platform{OS: platform.Windows, Arch: platform.Amd64}, version: "4.0", want: "Godot_v4.0-stable_win64.exe"},

		// v4.0-stable_mono
		{platform: platform.Platform{OS: platform.Windows, Arch: platform.I386}, version: "4.0-stable_mono", want: "Godot_v4.0-stable_mono_win32.exe"},
		{platform: platform.Platform{OS: platform.Windows, Arch: platform.Amd64}, version: "4.0-stable_mono", want: "Godot_v4.0-stable_mono_win64.exe"},

		// v5.0-rc4
		{platform: platform.Platform{OS: platform.Windows, Arch: platform.I386}, version: "5.0-rc4", want: "Godot_v5.0-rc4_win32.exe"},
		{platform: platform.Platform{OS: platform.Windows, Arch: platform.Amd64}, version: "5.0-rc4", want: "Godot_v5.0-rc4_win64.exe"},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d-%v-%s", i, tc.platform, tc.version), func(t *testing.T) {
			var v version.Version
			if tc.version != "" {
				v = version.MustParse(tc.version)
			}

			got, err := Executable{tc.platform, v}.Name()

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
		{s: "Godot_v1.0-stable_invalid", want: Executable{}, err: platform.ErrUnrecognizedPlatform},

		// Valid inputs
		// Linux
		{s: "Godot_v1.0-stable_x11.32", want: Executable{linux32(), v1}},
		{s: "Godot_v2.0-beta10_x11.64", want: Executable{linux64(), v2}},
		{s: "Godot_v3.0.4-alpha1_x11.32", want: Executable{linux32(), v3}},
		{s: "Godot_v4.0.11-dev.20230101_x11.64", want: Executable{linux64(), v4}},
		{s: "Godot_v4.0-stable_mono_linux_x86_64", want: Executable{linux64(), v4Mono}},

		// Darwin
		{s: "Godot_v1.0-stable_osx.fat", want: Executable{macOSUniversal(), v1}},
		{s: "Godot_v2.0-beta10_osx.64", want: Executable{macOSX86_64(), v2}},
		{s: "Godot_v3.0.4-alpha1_osx.universal", want: Executable{macOSUniversal(), v3}},
		{s: "Godot_v4.0.11-dev.20230101_macos.universal", want: Executable{macOSUniversal(), v4}},
		{s: "Godot_v4.0-stable_mono_macos.universal", want: Executable{macOSUniversal(), v4Mono}},

		// Windows
		{s: "Godot_v1.0-stable_win32", want: Executable{windows32(), v1}},
		{s: "Godot_v2.0-beta10_win64", want: Executable{windows64(), v2}},
		{s: "Godot_v3.0.4-alpha1_win32", want: Executable{windows32(), v3}},
		{s: "Godot_v4.0.11-dev.20230101_win64", want: Executable{windows64(), v4}},
		{s: "Godot_v4.0-stable_mono_win64", want: Executable{windows64(), v4Mono}},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d-%s", i, tc.s), func(t *testing.T) {
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

/* ---------------------- Functions: Platform Constants --------------------- */

// Returns a 'Platform' struct for 32-bit 'Linux'.
func linux32() platform.Platform {
	return platform.Platform{Arch: platform.I386, OS: platform.Linux}
}

// Returns a 'Platform' struct for 64-bit 'Linux'.
func linux64() platform.Platform {
	return platform.Platform{Arch: platform.Amd64, OS: platform.Linux}
}

// Returns a 'Platform' struct for 64-bit ('x86') 'MacOS'.
func macOSX86_64() platform.Platform {
	return platform.Platform{Arch: platform.Amd64, OS: platform.MacOS}
}

// Returns a 'Platform' struct for 64-bit ('ARM') 'MacOS'.
func macOSArm64() platform.Platform {
	return platform.Platform{Arch: platform.Arm64, OS: platform.MacOS}
}

// Returns a 'Platform' struct for a "fat" binary on 'MacOS'.
func macOSUniversal() platform.Platform {
	return platform.Platform{Arch: platform.Universal, OS: platform.MacOS}
}

// Returns a 'Platform' struct for 32-bit 'Windows'.
func windows32() platform.Platform {
	return platform.Platform{Arch: platform.I386, OS: platform.Windows}
}

// Returns a 'Platform' struct for 64-bit 'Windows'.
func windows64() platform.Platform {
	return platform.Platform{Arch: platform.Amd64, OS: platform.Windows}
}
