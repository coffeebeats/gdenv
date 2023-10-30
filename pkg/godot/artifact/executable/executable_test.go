package executable

import (
	"fmt"
	"testing"

	"github.com/coffeebeats/gdenv/pkg/godot/platform"
	"github.com/coffeebeats/gdenv/pkg/godot/version"
)

/* ------------------------------- Test: Name ------------------------------- */

func TestName(t *testing.T) {
	tests := []struct {
		platform platform.Platform
		version  string
		want     string
		err      error
	}{
		// Invalid inputs
		{platform: platform.Platform{}, want: ""},
		{platform: platform.Platform{OS: platform.Linux}, want: ""},
		{platform: linux64(), want: ""},

		{platform: platform.Platform{OS: platform.Linux, Arch: platform.Universal}, version: "4.0", want: ""},

		{platform: linux64(), version: "2.0", want: ""},
		{platform: macOSX86_64(), version: "2.0", want: ""},
		{platform: windows64(), version: "2.0", want: ""},

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
		{platform: windows32(), version: "3.6-beta1", want: "Godot_v3.6-beta1_win32.exe"},
		{platform: windows64(), version: "3.6-beta1", want: "Godot_v3.6-beta1_win64.exe"},

		// v4.0
		{platform: windows32(), version: "4.0", want: "Godot_v4.0-stable_win32.exe"},
		{platform: windows64(), version: "4.0", want: "Godot_v4.0-stable_win64.exe"},

		// v4.0-stable_mono
		{platform: windows32(), version: "4.0-stable_mono", want: "Godot_v4.0-stable_mono_win32.exe"},
		{platform: windows64(), version: "4.0-stable_mono", want: "Godot_v4.0-stable_mono_win64.exe"},

		// v5.0-rc4
		{platform: windows32(), version: "5.0-rc4", want: "Godot_v5.0-rc4_win32.exe"},
		{platform: windows64(), version: "5.0-rc4", want: "Godot_v5.0-rc4_win64.exe"},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d-%v-%s", i, tc.platform, tc.version), func(t *testing.T) {
			var v version.Version
			if tc.version != "" {
				v = version.MustParse(tc.version)
			}

			if got := (Executable{v, tc.platform}).Name(); got != tc.want {
				t.Errorf("output: got %#v, want %#v", got, tc.want)
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
