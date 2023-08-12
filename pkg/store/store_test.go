package store

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/coffeebeats/gdenv/internal/godot"
)

/* ------------------------------- Test: Init ------------------------------- */

func TestInit(t *testing.T) {
	tmp := t.TempDir()

	err := Init(tmp)
	if err != nil {
		t.Fatalf("err: %v", err)
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
			t.Fatalf("output: missing directory %s", d)
		}
	}

	for _, d := range []string{storeFileLayout} {
		if isDir, ok := got[d]; !ok || isDir {
			t.Fatalf("output: missing file %s", d)
		}
	}
}

/* -------------------------------- Test: Add ------------------------------- */

func TestAdd(t *testing.T) {
	tests := []struct {
		version godot.Version
		err     error
	}{
		{godot.Version{}, nil},
		{godot.Version{}.Canonical(), nil},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprint(tc.version), func(t *testing.T) {
			tmp := t.TempDir()
			store, tool := filepath.Join(tmp, "store"), filepath.Join(tmp, "tool")

			err := Init(store)
			if err != nil {
				t.Fatalf("err: %v", err)
			}

			// Create the tool to be moved.
			if err := os.WriteFile(tool, []byte(""), os.ModePerm); err != nil {
				t.Fatalf("test setup: %v", err)
			}

			// Invoke the 'Add' function.
			if err := Add(store, tool, tc.version); !errors.Is(err, tc.err) {
				t.Fatalf("err: got %#v, want %#v", err, tc.err)
			}

			// Verify the tool exists.
			name, err := godot.ExecutableName(tc.version.Canonical())
			if err != nil {
				t.Fatalf("test setup: %v", err)
			}

			toolWant := filepath.Join(store, storeDirGodot, tc.version.Canonical().String(), name)
			info, err := os.Stat(toolWant)
			if err != nil {
				t.Fatalf("output: %s", err)
			}

			if !info.Mode().IsRegular() {
				t.Fatalf("output is not a file: %s", toolWant)
			}
		})
	}
}

/* ------------------------------ Test: Remove ------------------------------ */

func TestRemove(t *testing.T) {
	tmp, version := t.TempDir(), godot.Version{}

	err := Init(tmp)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Create the tool to be moved.
	name, err := godot.ExecutableName(version.Canonical())
	if err != nil {
		t.Fatalf("test setup: %v", err)
	}

	toolWant := filepath.Join(tmp, storeDirGodot, version.Canonical().String(), name)
	if err := os.MkdirAll(filepath.Dir(toolWant), os.ModePerm); err != nil {
		t.Fatalf("test setup: %v", err)
	}
	if err := os.WriteFile(toolWant, []byte(""), os.ModePerm); err != nil {
		t.Fatalf("test setup: %v", err)
	}

	// Invoke the 'Remove' function.
	if err := Remove(tmp, version); !errors.Is(err, nil) {
		t.Fatalf("err: got %#v, want %#v", err, nil)
	}

	// Verify the tool is removed, along with the parent directory.
	info, err := os.Stat(filepath.Dir(toolWant))
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("output: %s", err)
	}

	if info != nil && info.Mode().IsDir() {
		t.Fatalf("output is not removed: %s", toolWant)
	}
}