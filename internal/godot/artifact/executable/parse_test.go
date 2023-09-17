package executable

// import (
// 	"errors"
// 	"testing"
// )

// /* -------------------------- Test: ParseExecutable ------------------------- */

// func TestParseExecutable(t *testing.T) {
// 	var (
// 		v1     = Version{major: 1}
// 		v2     = Version{major: 2, label: "beta10"}
// 		v3     = Version{3, 0, 4, "alpha1"}
// 		v4     = Version{4, 0, 11, "dev.20230101"}
// 		v4Mono = Version{4, 0, 0, "stable_mono"}
// 	)

// 	tests := []struct {
// 		s    string
// 		want Executable
// 		err  error
// 	}{
// 		// Invalid inputs
// 		{s: "", want: Executable{}, err: ErrMissingName},
// 		{s: "Godot-v1.0-stable-x11.32", want: Executable{}, err: ErrInvalidName},
// 		{s: "Godot-v1.0_stable-x11.32", want: Executable{}, err: ErrInvalidName},
// 		{s: "Godot_invalid_x11.32", want: Executable{}, err: ErrInvalidVersion},
// 		{s: "Godot_v1.0-stable_invalid", want: Executable{}, err: ErrUnrecognizedPlatform},

// 		// Valid inputs
// 		// Linux
// 		{s: "Godot_v1.0-stable_x11.32", want: Executable{Platform{i386, linux}, v1}, err: nil},
// 		{s: "Godot_v2.0-beta10_x11.64", want: Executable{Platform{amd64, linux}, v2}, err: nil},
// 		{s: "Godot_v3.0.4-alpha1_x11.32", want: Executable{Platform{i386, linux}, v3}, err: nil},
// 		{s: "Godot_v4.0.11-dev.20230101_x11.64", want: Executable{Platform{amd64, linux}, v4}, err: nil},
// 		{s: "Godot_v4.0-stable_mono_linux_x86_64", want: Executable{Platform{amd64, linux}, v4Mono}, err: nil},

// 		// Darwin
// 		{s: "Godot_v1.0-stable_osx.fat", want: Executable{Platform{universal, macOS}, v1}, err: nil},
// 		{s: "Godot_v2.0-beta10_osx.64", want: Executable{Platform{amd64, macOS}, v2}, err: nil},
// 		{s: "Godot_v3.0.4-alpha1_osx.universal", want: Executable{Platform{universal, macOS}, v3}, err: nil},
// 		{s: "Godot_v4.0.11-dev.20230101_macos.universal", want: Executable{Platform{universal, macOS}, v4}, err: nil},
// 		{s: "Godot_v4.0-stable_mono_macos.universal", want: Executable{Platform{universal, macOS}, v4Mono}, err: nil},

// 		// Windows
// 		{s: "Godot_v1.0-stable_win32", want: Executable{Platform{i386, windows}, v1}, err: nil},
// 		{s: "Godot_v2.0-beta10_win64", want: Executable{Platform{amd64, windows}, v2}, err: nil},
// 		{s: "Godot_v3.0.4-alpha1_win32", want: Executable{Platform{i386, windows}, v3}, err: nil},
// 		{s: "Godot_v4.0.11-dev.20230101_win64", want: Executable{Platform{amd64, windows}, v4}, err: nil},
// 		{s: "Godot_v4.0-stable_mono_win64", want: Executable{Platform{amd64, windows}, v4Mono}, err: nil},
// 	}

// 	for _, tc := range tests {
// 		t.Run(tc.s, func(t *testing.T) {
// 			got, err := ParseExecutable(tc.s)

// 			if !errors.Is(err, tc.err) {
// 				t.Fatalf("err: got %#v, want %#v", err, tc.err)
// 			}
// 			if got != tc.want {
// 				t.Fatalf("output: got %#v, want %#v", got, tc.want)
// 			}
// 		})
// 	}

// }
