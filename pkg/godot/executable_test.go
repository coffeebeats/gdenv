package godot

import (
	"errors"
	"fmt"
	"testing"
)

/* -------------------------- Test: ExecutableName -------------------------- */

func TestExecutableName(t *testing.T) {
	tests := []struct {
		p    Platform
		v    Version
		want string
		err  error
	}{
		// Invalid inputs
		{Platform{}, Version{}, "", ErrMissingOSInput},
		{Platform{os: linux}, Version{}, "", ErrMissingArchInput},
		{Platform{os: linux, arch: Amd64}, Version{}, "", ErrUnsupportedVersion},

		{Platform{os: linux, arch: Amd64}, Version{major: 2}, "", ErrUnsupportedVersion},
		{Platform{os: macOS, arch: Amd64}, Version{major: 2}, "", ErrUnsupportedVersion},
		{Platform{os: windows, arch: Amd64}, Version{major: 2}, "", ErrUnsupportedVersion},

		// Valid inputs - Linux

		// v3.0
		{Platform{os: linux, arch: I386}, Version{major: 3}, "Godot_v3.0-stable_x11.32", nil},
		{Platform{os: linux, arch: Amd64}, Version{major: 3}, "Godot_v3.0-stable_x11.64", nil},
		{Platform{os: linux, arch: Arm64}, Version{major: 3}, "", ErrUnsupportedArchInput},
		{Platform{os: linux, arch: Universal}, Version{major: 3}, "", ErrUnsupportedArchInput},

		// v3.6-beta2
		{Platform{os: linux, arch: I386}, Version{major: 3, minor: 6, label: "beta2"}, "Godot_v3.6-beta2_x11.32", nil},
		{Platform{os: linux, arch: Amd64}, Version{major: 3, minor: 6, label: "beta2"}, "Godot_v3.6-beta2_x11.64", nil},
		{Platform{os: linux, arch: Arm64}, Version{major: 3, minor: 6, label: "beta2"}, "", ErrUnsupportedArchInput},
		{Platform{os: linux, arch: Universal}, Version{major: 3, minor: 6, label: "beta2"}, "", ErrUnsupportedArchInput},

		// v4.0-rc4
		{Platform{os: linux, arch: I386}, Version{major: 4, label: "rc4"}, "Godot_v4.0-rc4_linux.x86_32", nil},
		{Platform{os: linux, arch: Amd64}, Version{major: 4, label: "rc4"}, "Godot_v4.0-rc4_linux.x86_64", nil},
		{Platform{os: linux, arch: Arm64}, Version{major: 4, label: "rc4"}, "", ErrUnsupportedArchInput},
		{Platform{os: linux, arch: Universal}, Version{major: 4, label: "rc4"}, "", ErrUnsupportedArchInput},

		// Valid inputs - MacOS

		// v3.0
		{Platform{os: macOS, arch: I386}, Version{major: 3}, "Godot_v3.0-stable_osx.fat", nil},
		{Platform{os: macOS, arch: Amd64}, Version{major: 3}, "Godot_v3.0-stable_osx.fat", nil},
		{Platform{os: macOS, arch: Arm64}, Version{major: 3}, "", ErrUnsupportedArchInput},
		{Platform{os: macOS, arch: Universal}, Version{major: 3}, "", ErrUnsupportedArchInput},

		// v3.6-beta2
		{Platform{os: macOS, arch: Amd64}, Version{major: 3, minor: 6, label: "beta2"}, "Godot_v3.6-beta2_osx.universal", nil},
		{Platform{os: macOS, arch: Arm64}, Version{major: 3, minor: 6, label: "beta2"}, "Godot_v3.6-beta2_osx.universal", nil},
		{Platform{os: macOS, arch: I386}, Version{major: 3, minor: 6, label: "beta2"}, "", ErrUnsupportedArchInput},
		{Platform{os: macOS, arch: Universal}, Version{major: 3, minor: 6, label: "beta2"}, "", ErrUnsupportedArchInput},

		// v4.0-rc4
		{Platform{os: macOS, arch: Amd64}, Version{major: 4, label: "rc4"}, "Godot_v4.0-rc4_macos.universal", nil},
		{Platform{os: macOS, arch: Arm64}, Version{major: 4, label: "rc4"}, "Godot_v4.0-rc4_macos.universal", nil},
		{Platform{os: macOS, arch: I386}, Version{major: 4, label: "rc4"}, "", ErrUnsupportedArchInput},
		{Platform{os: macOS, arch: Universal}, Version{major: 4, label: "rc4"}, "", ErrUnsupportedArchInput},

		// Valid inputs - Windows

		// v3.0
		{Platform{os: windows, arch: I386}, Version{major: 3}, "Godot_v3.0-stable_win32.exe", nil},
		{Platform{os: windows, arch: Amd64}, Version{major: 3}, "Godot_v3.0-stable_win64.exe", nil},
		{Platform{os: windows, arch: Arm64}, Version{major: 3}, "", ErrUnsupportedArchInput},
		{Platform{os: windows, arch: Universal}, Version{major: 3}, "", ErrUnsupportedArchInput},

		// v3.6-beta2
		{Platform{os: windows, arch: I386}, Version{major: 3, minor: 6, label: "beta2"}, "Godot_v3.6-beta2_win32.exe", nil},
		{Platform{os: windows, arch: Amd64}, Version{major: 3, minor: 6, label: "beta2"}, "Godot_v3.6-beta2_win64.exe", nil},
		{Platform{os: windows, arch: Arm64}, Version{major: 3, minor: 6, label: "beta2"}, "", ErrUnsupportedArchInput},
		{Platform{os: windows, arch: Universal}, Version{major: 3, minor: 6, label: "beta2"}, "", ErrUnsupportedArchInput},

		// v4.0-rc4
		{Platform{os: windows, arch: I386}, Version{major: 4, label: "rc4"}, "Godot_v4.0-rc4_win32.exe", nil},
		{Platform{os: windows, arch: Amd64}, Version{major: 4, label: "rc4"}, "Godot_v4.0-rc4_win64.exe", nil},
		{Platform{os: windows, arch: Arm64}, Version{major: 4, label: "rc4"}, "", ErrUnsupportedArchInput},
		{Platform{os: windows, arch: Universal}, Version{major: 4, label: "rc4"}, "", ErrUnsupportedArchInput},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got, err := ExecutableName(tc.p, tc.v)

			if !errors.Is(err, tc.err) {
				t.Fatalf("err: got %#v, want %#v", err, tc.err)
			}
			if got != tc.want {
				t.Fatalf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}
