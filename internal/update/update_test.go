package update

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/coffeebeats/gdenv/pkg/store"
)

/* -------------------------------------------------------------------------- */
/*                             Function: TestCheck                            */
/* -------------------------------------------------------------------------- */

func TestCheck(t *testing.T) {
	// Given: A version which cannot be parsed as a semantic version.
	current := "not-a-version"

	// When: The latest release is checked against it.
	_, err := Check(context.Background(), current)

	// Then: The version is rejected rather than reported as up-to-date.
	//
	// NOTE: This returns before any request is made, so no network is involved.
	if !errors.Is(err, ErrInvalidVersion) {
		t.Fatalf("err: got %v, want %v", err, ErrInvalidVersion)
	}
}

/* -------------------------------------------------------------------------- */
/*                         Function: TestManagedFrom                          */
/* -------------------------------------------------------------------------- */

func TestManagedFrom(t *testing.T) {
	t.Run("accepts a binary within the store's bin directory", func(t *testing.T) {
		home := t.TempDir()
		binDir := filepath.Join(home, "bin")

		if err := os.Mkdir(binDir, 0o755); err != nil {
			t.Fatal(err)
		}

		t.Setenv("GDENV_HOME", home)

		exe := filepath.Join(binDir, "gdenv")
		writeFile(t, exe, "binary")

		got, err := managedBinDir(exe)
		if err != nil {
			t.Fatal(err)
		}

		if !samePath(got, binDir) {
			t.Fatalf("got %q, want %q", got, binDir)
		}
	})

	t.Run("rejects a binary outside the store", func(t *testing.T) {
		home, elsewhere := t.TempDir(), t.TempDir()

		t.Setenv("GDENV_HOME", home)

		exe := filepath.Join(elsewhere, "gdenv")
		writeFile(t, exe, "binary")

		_, err := managedBinDir(exe)
		if !errors.Is(err, ErrUnmanaged) {
			t.Fatalf("err: got %v, want %v", err, ErrUnmanaged)
		}

		// The message must name both paths so the user can act on it.
		for _, want := range []string{exe, "bin"} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("error %q does not mention %q", err, want)
			}
		}
	})

	t.Run("reports a missing GDENV_HOME", func(t *testing.T) {
		t.Setenv("GDENV_HOME", "")

		if _, err := managedBinDir(filepath.Join(t.TempDir(), "gdenv")); !errors.Is(err, store.ErrMissingEnvVar) {
			t.Fatalf("err: got %v, want %v", err, store.ErrMissingEnvVar)
		}
	})

	t.Run("accepts a symlinked GDENV_HOME", func(t *testing.T) {
		if runtime.GOOS == osWindows {
			t.Skip("creating symlinks on Windows requires elevated privileges")
		}

		real := t.TempDir()

		binDir := filepath.Join(real, "bin")
		if err := os.Mkdir(binDir, 0o755); err != nil {
			t.Fatal(err)
		}

		exe := filepath.Join(binDir, "gdenv")
		writeFile(t, exe, "binary")

		// Point 'GDENV_HOME' at a symlink to the real directory; a naive path
		// comparison would reject this valid installation.
		link := filepath.Join(t.TempDir(), "home")
		if err := os.Symlink(real, link); err != nil {
			t.Fatal(err)
		}

		t.Setenv("GDENV_HOME", link)

		if _, err := managedBinDir(exe); err != nil {
			t.Fatalf("expected a symlinked home to be accepted, got %v", err)
		}
	})

	t.Run("accepts a case-differing path on Windows", func(t *testing.T) {
		if runtime.GOOS != osWindows {
			t.Skip("path comparison is case-sensitive off Windows")
		}

		home := t.TempDir()

		binDir := filepath.Join(home, "bin")
		if err := os.Mkdir(binDir, 0o755); err != nil {
			t.Fatal(err)
		}

		exe := filepath.Join(binDir, "gdenv.exe")
		writeFile(t, exe, "binary")

		t.Setenv("GDENV_HOME", strings.ToUpper(home))

		if _, err := managedBinDir(exe); err != nil {
			t.Fatalf("expected a case-differing path to be accepted, got %v", err)
		}
	})
}

/* -------------------------------------------------------------------------- */
/*                        Function: TestVerifyWritable                        */
/* -------------------------------------------------------------------------- */

func TestVerifyWritable(t *testing.T) {
	dir := t.TempDir()

	if err := verifyWritable(dir); err != nil {
		t.Fatalf("expected a temporary directory to be writable, got %v", err)
	}

	// The probe file is not swept by 'cleanStale', so it must not be left behind.
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	if len(entries) != 0 {
		t.Fatalf("entries: got %v, want none", entries)
	}

	if err := verifyWritable(filepath.Join(t.TempDir(), "does-not-exist")); !errors.Is(err, ErrNotWritable) {
		t.Fatalf("err: got %v, want %v", err, ErrNotWritable)
	}
}

/* -------------------------------------------------------------------------- */
/*                         Function: TestAcquireLock                          */
/* -------------------------------------------------------------------------- */

func TestAcquireLock(t *testing.T) {
	t.Run("is exclusive while held and reusable once released", func(t *testing.T) {
		dir := t.TempDir()

		unlock, err := acquireLock(dir)
		if err != nil {
			t.Fatal(err)
		}

		if _, err := acquireLock(dir); !errors.Is(err, ErrLockHeld) {
			t.Fatalf("err: got %v, want %v", err, ErrLockHeld)
		}

		unlock()

		unlock, err = acquireLock(dir)
		if err != nil {
			t.Fatalf("expected the lock to be reacquirable, got %v", err)
		}

		unlock()
	})

	t.Run("reclaims an abandoned lock", func(t *testing.T) {
		dir := t.TempDir()

		path := filepath.Join(dir, lockName)
		writeFile(t, path, "")

		// Backdate the lock past the staleness threshold, as though an earlier
		// update had been killed partway through.
		old := time.Now().Add(-lockStaleAfter * 2)
		if err := os.Chtimes(path, old, old); err != nil {
			t.Fatal(err)
		}

		unlock, err := acquireLock(dir)
		if err != nil {
			t.Fatalf("expected a stale lock to be reclaimed, got %v", err)
		}

		unlock()
	})
}
