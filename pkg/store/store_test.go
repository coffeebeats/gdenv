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
		err         error
	}{
		{"linux", "amd64", "4.0", nil},
		{"linux", "amd64", "4.0-alpha1", nil},

		{"macos", "amd64", "4.0", nil},
		{"macos", "amd64", "4.0-alpha1", nil},

		{"windows", "i386", "4.0", nil},
		{"windows", "i386", "4.0-alpha1", nil},
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
			if err := os.WriteFile(tool, []byte(""), os.ModePerm); err != nil {
				t.Fatalf("test setup: %v", err)
			}

			// Invoke the 'Add' function.
			if err := Add(store, tool, ex); !errors.Is(err, tc.err) {
				t.Errorf("err: got %#v, want %#v", err, tc.err)
			}

			// Verify the tool exists.
			toolWant := filepath.Join(store, storeDirGodot, v.String(), ex.Name())
			info, err := os.Stat(toolWant)
			if err != nil {
				t.Errorf("output: %s", err)
			}

			if !info.Mode().IsRegular() {
				t.Errorf("output is not a file: %s", toolWant)
			}
		})
	}
}

/* ------------------------------ Test: Remove ------------------------------ */

func TestRemove(t *testing.T) {
	tests := []struct {
		os, arch, v string
		err         error
	}{
		{"linux", "amd64", "4.0", nil},
		{"linux", "amd64", "4.0-alpha1", nil},

		{"macos", "amd64", "4.0", nil},
		{"macos", "amd64", "4.0-alpha1", nil},

		{"windows", "i386", "4.0", nil},
		{"windows", "i386", "4.0-alpha1", nil},
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
			toolWant := filepath.Join(tmp, storeDirGodot, v.String(), ex.Name())
			if err := os.MkdirAll(filepath.Dir(toolWant), os.ModePerm); err != nil {
				t.Fatalf("test setup: %v", err)
			}
			if err := os.WriteFile(toolWant, []byte(""), os.ModePerm); err != nil {
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
