package store

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/coffeebeats/gdenv/internal/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/internal/godot/platform"
	"github.com/coffeebeats/gdenv/internal/godot/version"
)

const (
	modeMockTool = 0600 // rw-------
	modeTestDir  = 0700 // rwx------
)

/* ------------------------------- Test: Init ------------------------------- */

func TestInit(t *testing.T) {
	tmp := t.TempDir()

	err := Init(tmp)
	if err != nil {
		t.Errorf("err: %v", err)
	}

	files, err := os.ReadDir(tmp)
	if err != nil {
		t.Fatalf("test setup: %v", err)
	}

	got := make(map[string]bool, len(files))
	for _, f := range files {
		got[f.Name()] = f.IsDir()
	}

	for _, d := range []string{storeDirBin, storeDirGodot} {
		if isDir, ok := got[d]; !ok || !isDir {
			t.Errorf("output: missing directory %s", d)
		}
	}

	for _, d := range []string{storeFileLayout} {
		if isDir, ok := got[d]; !ok || isDir {
			t.Errorf("output: missing file %s", d)
		}
	}
}

/* -------------------------------- Test: Add ------------------------------- */

func TestAdd(t *testing.T) {
	tests := []struct {
		os, arch, v string

		want string
	}{
		{os: "linux", arch: "amd64", v: "4.0", want: "v4.0-stable/Godot_v4.0-stable_linux.x86_64"},
		{os: "linux", arch: "amd64", v: "4.0-alpha1", want: "v4.0-alpha1/Godot_v4.0-alpha1_linux.64"},

		{os: "macos", arch: "amd64", v: "4.0", want: "v4.0-stable/Godot.app"},
		{os: "macos", arch: "amd64", v: "4.0-alpha1", want: "v4.0-alpha1/Godot.app"},

		{os: "windows", arch: "i386", v: "4.0", want: "v4.0-stable/Godot_v4.0-stable_win32.exe"},
		{os: "windows", arch: "i386", v: "4.0-alpha1", want: "v4.0-alpha1/Godot_v4.0-alpha1_win32.exe"},
	}

	for _, tc := range tests {
		t.Run(tc.v, func(t *testing.T) {
			tmp := t.TempDir()
			store, tool := filepath.Join(tmp, "store"), filepath.Join(tmp, "tool")

			err := Init(store)
			if err != nil {
				t.Errorf("err: %v", err)
			}

			// Define the 'Version' for the test.
			v := version.MustParse(tc.v)

			// Define the 'Platform' for the test.
			p := platform.Platform{
				Arch: platform.MustParseArch(tc.arch),
				OS:   platform.MustParseOS(tc.os),
			}

			// Define the 'Executable' for the test.
			ex := executable.New(v, p)

			// Create the tool to be moved.
			if err := os.WriteFile(tool, []byte(""), modeMockTool); err != nil {
				t.Fatalf("test setup: %v", err)
			}

			// Invoke the 'Add' function.
			if err := Add(store, tool, ex); !errors.Is(err, nil) {
				t.Errorf("err: got %#v, want %#v", err, nil)
			}

			// Verify the tool exists.
			toolWant := filepath.Join(store, storeDirGodot, tc.want)
			if _, err = os.Stat(toolWant); err != nil {
				t.Errorf("err: %v", err)
			}
		})
	}
}

/* ------------------------------ Test: Remove ------------------------------ */

func TestRemove(t *testing.T) {
	tests := []struct {
		os, arch, v string

		want string
	}{
		{os: "linux", arch: "amd64", v: "4.0", want: "v4.0-stable/Godot_v4.0-stable_linux.x86_64"},
		{os: "linux", arch: "amd64", v: "4.0-alpha1", want: "v4.0-alpha1/Godot_v4.0-alpha1_linux.64"},

		{os: "macos", arch: "amd64", v: "4.0", want: "v4.0-stable/Godot.app"},
		{os: "macos", arch: "amd64", v: "4.0-alpha1", want: "v4.0-alpha1/Godot.app"},

		{os: "windows", arch: "i386", v: "4.0", want: "v4.0-stable/Godot_v4.0-stable_win32.exe"},
		{os: "windows", arch: "i386", v: "4.0-alpha1", want: "v4.0-alpha1/Godot_v4.0-alpha1_win32.exe"},
	}

	for _, tc := range tests {

		t.Run(tc.v, func(t *testing.T) {
			tmp := t.TempDir()

			err := Init(tmp)
			if err != nil {
				t.Errorf("err: %v", err)
			}

			// Define the 'Version' for the test.
			v := version.MustParse(tc.v)

			// Define the 'Platform' for the test.
			p := platform.Platform{
				Arch: platform.MustParseArch(tc.arch),
				OS:   platform.MustParseOS(tc.os),
			}

			// Define the 'Executable' for the test.
			ex := executable.New(v, p)

			// Create the tool to be moved.
			toolWant := filepath.Join(tmp, storeDirGodot, tc.want)
			if err := os.MkdirAll(filepath.Dir(toolWant), modeTestDir); err != nil {
				t.Fatalf("test setup: %v", err)
			}
			if err := os.WriteFile(toolWant, []byte(""), modeMockTool); err != nil {
				t.Fatalf("test setup: %v", err)
			}

			// Invoke the 'Remove' function.
			if err := Remove(tmp, ex); !errors.Is(err, nil) {
				t.Errorf("err: got %#v, want %#v", err, nil)
			}

			// Verify the tool is removed, along with the parent directory.
			info, err := os.Stat(filepath.Dir(toolWant))
			if !errors.Is(err, fs.ErrNotExist) {
				t.Errorf("output: %s", err)
			}

			if info != nil && info.Mode().IsDir() {
				t.Errorf("output is not removed: %s", toolWant)
			}
		})
	}

}
