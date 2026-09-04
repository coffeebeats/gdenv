package update

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

var errForcedRename = errors.New("forced rename failure")

/* -------------------------------------------------------------------------- */
/*                        Function: TestReplaceBinary                         */
/* -------------------------------------------------------------------------- */

func TestReplaceBinary(t *testing.T) {
	t.Run("replaces the target and preserves the original", func(t *testing.T) {
		dir := t.TempDir()
		src, target := filepath.Join(dir, "new"), filepath.Join(dir, "gdenv")

		writeFile(t, src, "new-binary")
		writeFile(t, target, "old-binary")

		restore, err := replaceBinary(src, target)
		if err != nil {
			t.Fatal(err)
		}

		if got := readFile(t, target); got != "new-binary" {
			t.Fatalf("target: got %q, want %q", got, "new-binary")
		}

		// The original must remain recoverable until the update is committed.
		if got := readFile(t, target+extensionStale); got != "old-binary" {
			t.Fatalf("stale: got %q, want %q", got, "old-binary")
		}

		if err := restore(); err != nil {
			t.Fatal(err)
		}

		if got := readFile(t, target); got != "old-binary" {
			t.Fatalf("after restore: got %q, want %q", got, "old-binary")
		}
	})

	t.Run("succeeds when the target does not exist", func(t *testing.T) {
		dir := t.TempDir()
		src, target := filepath.Join(dir, "new"), filepath.Join(dir, "godot")

		writeFile(t, src, "new-binary")

		restore, err := replaceBinary(src, target)
		if err != nil {
			t.Fatal(err)
		}

		if got := readFile(t, target); got != "new-binary" {
			t.Fatalf("target: got %q, want %q", got, "new-binary")
		}

		if err := restore(); err != nil {
			t.Fatal(err)
		}

		if _, err := os.Stat(target); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("expected target to be removed, got %v", err)
		}
	})

	t.Run("restores the original when the replacement fails", func(t *testing.T) {
		dir := t.TempDir()
		src, target := filepath.Join(dir, "new"), filepath.Join(dir, "gdenv")

		writeFile(t, src, "new-binary")
		writeFile(t, target, "old-binary")

		// Fail only the second rename, which moves the new binary into place.
		withRenameFailure(t, "new")

		if _, err := replaceBinary(src, target); !errors.Is(err, errForcedRename) {
			t.Fatalf("err: got %v, want %v", err, errForcedRename)
		}

		// The original must be back in place, not left as a '.old' file.
		if got := readFile(t, target); got != "old-binary" {
			t.Fatalf("target: got %q, want %q", got, "old-binary")
		}
	})

	t.Run("removes a leftover stale file first", func(t *testing.T) {
		dir := t.TempDir()
		src, target := filepath.Join(dir, "new"), filepath.Join(dir, "gdenv")

		writeFile(t, src, "new-binary")
		writeFile(t, target, "old-binary")
		writeFile(t, target+extensionStale, "ancient-binary")

		if _, err := replaceBinary(src, target); err != nil {
			t.Fatal(err)
		}

		if got := readFile(t, target+extensionStale); got != "old-binary" {
			t.Fatalf("stale: got %q, want %q", got, "old-binary")
		}
	})
}

/* -------------------------------------------------------------------------- */
/*                         Function: TestReplaceAll                           */
/* -------------------------------------------------------------------------- */

func TestReplaceAll(t *testing.T) {
	t.Run("replaces every binary and removes the superseded ones", func(t *testing.T) {
		binDir, srcDir := t.TempDir(), t.TempDir()
		names := []string{"godot", "gdenv"}

		found := make(map[string]string, len(names))

		for _, name := range names {
			writeFile(t, filepath.Join(binDir, name), "old-"+name)

			src := filepath.Join(srcDir, name)
			writeFile(t, src, "new-"+name)
			found[name] = src
		}

		if err := replaceAll(binDir, names, found); err != nil {
			t.Fatal(err)
		}

		for _, name := range names {
			if got := readFile(t, filepath.Join(binDir, name)); got != "new-"+name {
				t.Fatalf("%s: got %q, want %q", name, got, "new-"+name)
			}

			stale := filepath.Join(binDir, name+extensionStale)
			if _, err := os.Stat(stale); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("expected %s to be removed, got %v", stale, err)
			}
		}
	})

	t.Run("rolls back an earlier binary when a later one fails", func(t *testing.T) {
		binDir, srcDir := t.TempDir(), t.TempDir()

		// 'gdenv' is replaced last, mirroring the real ordering.
		names := []string{"godot", "gdenv"}

		found := make(map[string]string, len(names))

		for _, name := range names {
			writeFile(t, filepath.Join(binDir, name), "old-"+name)

			src := filepath.Join(srcDir, name)
			writeFile(t, src, "new-"+name)
			found[name] = src
		}

		// Fail while replacing 'gdenv', after 'godot' has already succeeded.
		withRenameFailure(t, "gdenv")

		if err := replaceAll(binDir, names, found); !errors.Is(err, errForcedRename) {
			t.Fatalf("err: got %v, want %v", err, errForcedRename)
		}

		// Both binaries must be back to their original contents.
		for _, name := range names {
			if got := readFile(t, filepath.Join(binDir, name)); got != "old-"+name {
				t.Fatalf("%s: got %q, want %q after rollback", name, got, "old-"+name)
			}
		}
	})
}

/* -------------------------------------------------------------------------- */
/*                          Function: TestCleanStale                          */
/* -------------------------------------------------------------------------- */

func TestCleanStale(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, filepath.Join(dir, "gdenv"), "keep")
	writeFile(t, filepath.Join(dir, "gdenv"+extensionStale), "sweep")
	writeFile(t, filepath.Join(dir, "godot.exe"+extensionStale), "sweep")

	staleDir := filepath.Join(dir, prefixTempDir+"123")
	if err := os.Mkdir(staleDir, 0o755); err != nil {
		t.Fatal(err)
	}

	writeFile(t, filepath.Join(staleDir, "leftover"), "sweep")

	cleanStale(dir)

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	if len(entries) != 1 || entries[0].Name() != "gdenv" {
		names := make([]string, 0, len(entries))
		for _, e := range entries {
			names = append(names, e.Name())
		}

		t.Fatalf("remaining: got %v, want [gdenv]", names)
	}
}

/* -------------------------------------------------------------------------- */
/*                              Test Helpers                                  */
/* -------------------------------------------------------------------------- */

// withRenameFailure replaces the package's rename implementation so that a
// failure can be simulated on any platform, and restores it afterwards.
//
// NOTE: This is what makes the Windows-specific path - where the running
// executable cannot be overwritten - reachable on Linux CI.
func withRenameFailure(t *testing.T, failOn string) {
	t.Helper()

	original := rename

	rename = func(oldpath, newpath string) error {
		if filepath.Base(newpath) == failOn || filepath.Base(oldpath) == failOn {
			return errForcedRename
		}

		return original(oldpath, newpath)
	}

	t.Cleanup(func() { rename = original })
}

func writeFile(t *testing.T, path, contents string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	return string(b)
}
