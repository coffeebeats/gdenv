package executable

// import (
// 	"errors"
// 	"fmt"
// 	"testing"

// 	"github.com/coffeebeats/gdenv/pkg/godot"
// )

// /* -------------------------- Test: Executable.Name ------------------------- */

// func TestExecutableName(t *testing.T) {
// 	var (
// 		v3     = godot.NewVersionWithLabel(3, 6, 0, "beta1")
// 		v4     = godot.NewVersion(4, 0, 0)
// 		v4Mono = godot.NewVersionWithLabel(4, 0, 0, "stable_mono")
// 		v5     = godot.NewVersionWithLabel(4, 0, 0, "rc4")
// 	)

// 	var (
// 		linuxAmd64     = godot.MustParsePlatform("linux.x86_64")
// 		linuxArm64     = godot.MustParsePlatform("linux.x86_64")
// 		linuxI386      = godot.MustParsePlatform("linux.x86_64")
// 		linuxUniversal = godot.MustParsePlatform("linux.x86_64")

// 		macOSAmd64     = godot.MustParsePlatform("osx.64")
// 		macOSArm64     = godot.MustParsePlatform("osx.arm64")
// 		macOSI386      = godot.MustParsePlatform("linux.x86_64")
// 		macOSUniversal = godot.MustParsePlatform("linux.x86_64")

// 		windowsAmd64     = godot.MustParsePlatform("linux.x86_64")
// 		windowsArm64     = godot.MustParsePlatform("linux.x86_64")
// 		windowsI386      = godot.MustParsePlatform("linux.x86_64")
// 		windowsUniversal = godot.MustParsePlatform("linux.x86_64")
// 	)

// 	tests := []struct {
// 		platform godot.Platform
// 		version  godot.Version
// 		want     string
// 	}{
// 		// Invalid inputs
// 		{godot.Platform{}, godot.Version{}, ""},

// 		// Valid inputs - Linux

// 		// v3.6-beta1
// 		{platform: linuxI386, version: v3, want: "Godot_v3.6-beta1_x11.32"},
// 		{platform: linuxAmd64, version: v3, want: "Godot_v3.6-beta1_x11.64"},
// 		{platform: linuxArm64, version: v3, want: ""},
// 		{platform: linuxUniversal, version: v3, want: ""},

// 		// v4.0
// 		{platform: linuxI386, version: v4, want: "Godot_v4.0-stable_linux.x86_32"},
// 		{platform: linuxAmd64, version: v4, want: "Godot_v4.0-stable_linux.x86_64"},
// 		{platform: linuxArm64, version: v4, want: ""},
// 		{platform: linuxUniversal, version: v4, want: ""},

// 		// v4.0-stable_mono
// 		{platform: linuxI386, version: v4Mono, want: "Godot_v4.0-stable_mono_linux_x86_32"},
// 		{platform: linuxAmd64, version: v4Mono, want: "Godot_v4.0-stable_mono_linux_x86_64"},
// 		{platform: linuxArm64, version: v4Mono, want: ""},
// 		{platform: linuxUniversal, version: v4Mono, want: ""},

// 		// v5.0-rc4
// 		{platform: linuxI386, version: v5, want: "Godot_v5.0-rc4_linux.x86_32"},
// 		{platform: linuxAmd64, version: v5, want: "Godot_v5.0-rc4_linux.x86_64"},
// 		{platform: linuxArm64, version: v5, want: ""},
// 		{platform: linuxUniversal, version: v5, want: ""},

// 		// Valid inputs - MacOS

// 		// v3.6-beta1
// 		{platform: macOSAmd64, version: v3, want: "Godot_v3.6-beta1_osx.universal"},
// 		{platform: macOSArm64, version: v3, want: "Godot_v3.6-beta1_osx.universal"},
// 		{platform: godot.Platform{os: macOS, arch: i386}, version: v3, want: ""},
// 		{platform: macOSUniversal, version: v3, want: ""},

// 		// v4.0
// 		{platform: macOSAmd64, version: v4, want: "Godot_v4.0-stable_macos.universal"},
// 		{platform: macOSArm64, version: v4, want: "Godot_v4.0-stable_macos.universal"},
// 		{platform: godot.Platform{os: macOS, arch: i386}, version: v4, want: ""},
// 		{platform: macOSUniversal, version: v4, want: ""},

// 		// v4.0-stable_mono
// 		{platform: macOSAmd64, version: v4Mono, want: "Godot_v4.0-stable_mono_macos.universal"},
// 		{platform: macOSArm64, version: v4Mono, want: "Godot_v4.0-stable_mono_macos.universal"},
// 		{platform: godot.Platform{os: macOS, arch: i386}, version: v4Mono, want: ""},
// 		{platform: macOSUniversal, version: v4Mono, want: ""},

// 		// v5.0-rc4
// 		{platform: macOSAmd64, version: v5, want: "Godot_v5.0-rc4_macos.universal"},
// 		{platform: macOSArm64, version: v5, want: "Godot_v5.0-rc4_macos.universal"},
// 		{platform: godot.Platform{os: macOS, arch: i386}, version: v5, want: ""},
// 		{platform: macOSUniversal, version: v5, want: ""},

// 		// Valid inputs - Windows

// 		// v3.6-beta1
// 		{platform: godot.Platform{os: windows, arch: i386}, version: v3, want: "Godot_v3.6-beta1_win32.exe"},
// 		{platform: godot.Platform{os: windows, arch: amd64}, version: v3, want: "Godot_v3.6-beta1_win64.exe"},
// 		{platform: godot.Platform{os: windows, arch: arm64}, version: v3, want: ""},
// 		{platform: godot.Platform{os: windows, arch: universal}, version: v3, want: ""},

// 		// v4.0
// 		{platform: godot.Platform{os: windows, arch: i386}, version: v4, want: "Godot_v4.0-stable_win32.exe"},
// 		{platform: godot.Platform{os: windows, arch: amd64}, version: v4, want: "Godot_v4.0-stable_win64.exe"},
// 		{platform: godot.Platform{os: windows, arch: arm64}, version: v4, want: ""},
// 		{platform: godot.Platform{os: windows, arch: universal}, version: v4, want: ""},

// 		// v4.0-stable_mono
// 		{platform: godot.Platform{os: windows, arch: i386}, version: v4Mono, want: "Godot_v4.0-stable_mono_win32.exe"},
// 		{platform: godot.Platform{os: windows, arch: amd64}, version: v4Mono, want: "Godot_v4.0-stable_mono_win64.exe"},
// 		{platform: godot.Platform{os: windows, arch: arm64}, version: v4Mono, want: ""},
// 		{platform: godot.Platform{os: windows, arch: universal}, version: v4Mono, want: ""},

// 		// v5.0-rc4
// 		{platform: godot.Platform{os: windows, arch: i386}, version: v5, want: "Godot_v5.0-rc4_win32.exe"},
// 		{platform: godot.Platform{os: windows, arch: amd64}, version: v5, want: "Godot_v5.0-rc4_win64.exe"},
// 		{platform: godot.Platform{os: windows, arch: arm64}, version: v5, want: ""},
// 		{platform: godot.Platform{os: windows, arch: universal}, version: v5, want: ""},
// 	}

// 	for _, tc := range tests {
// 		t.Run(fmt.Sprintf("%v-%s", tc.p, tc.v), func(t *testing.T) {
// 			got, err := Executable{tc.p, tc.v}.Name()

// 			if !errors.Is(err, tc.err) {
// 				t.Fatalf("err: got %#v, want %#v", err, tc.err)
// 			}
// 			if got != tc.want {
// 				t.Fatalf("output: got %#v, want %#v", got, tc.want)
// 			}
// 		})
// 	}
// }
