// Package update implements self-updating for the 'gdenv' binary.
//
// The mechanism mirrors what 'scripts/install.sh' and 'scripts/install.ps1' do
// - download the release archive published for the host, verify it against the
// published checksums, and unpack it into '$GDENV_HOME/bin' - with two
// differences: the target platform is taken from the running binary rather than
// detected, and the shell profile is left alone because installation has
// already happened.
package update

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/log"

	"github.com/coffeebeats/gdenv/internal/osutil"
	"github.com/coffeebeats/gdenv/pkg/store"
)

const (
	lockName       = ".update.lock"
	lockStaleAfter = 10 * time.Minute
)

var (
	ErrLockHeld    = errors.New("another update is already in progress")
	ErrNotWritable = errors.New("install directory is not writable")
	ErrUnmanaged   = errors.New("binary was not installed by gdenv's installer")
)

/* -------------------------------------------------------------------------- */
/*                               Struct: Status                               */
/* -------------------------------------------------------------------------- */

// Status describes the outcome of an update check.
type Status struct {
	Current string
	Latest  string
}

/* ---------------------------- Method: IsUpgrade --------------------------- */

// IsUpgrade reports whether the latest version is newer than the current one.
func (s Status) IsUpgrade() bool {
	return IsUpgrade(s.Current, s.Latest)
}

/* -------------------------------------------------------------------------- */
/*                               Function: Check                              */
/* -------------------------------------------------------------------------- */

// Check resolves the latest published version and compares it to the running
// one.
func Check(ctx context.Context, current string) (Status, error) {
	var status Status

	// NOTE: A version which cannot be parsed must be rejected rather than
	// compared; 'IsUpgrade' would report 'false' for it, which the caller cannot
	// distinguish from a genuinely up-to-date installation.
	current = Normalize(current)
	if !IsValidVersion(current) {
		return status, fmt.Errorf("%w: %s", ErrInvalidVersion, current)
	}

	latest, err := LatestVersion(ctx)
	if err != nil {
		return status, err
	}

	status.Current = current
	status.Latest = latest

	return status, nil
}

/* -------------------------------------------------------------------------- */
/*                           Function: ManagedBinDir                          */
/* -------------------------------------------------------------------------- */

// ManagedBinDir returns the directory holding the running 'gdenv' binary, after
// verifying that it is the directory managed by gdenv's installer.
//
// Refusing to touch anything else keeps 'gdenv' from overwriting a binary owned
// by a system package manager or produced by 'go install'.
func ManagedBinDir() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}

	return managedBinDir(exe)
}

/* --------------------------- Function: managedBinDir ---------------------- */

// managedBinDir implements 'ManagedBinDir' for an explicit executable path.
func managedBinDir(exe string) (string, error) {
	exe = resolveSymlinks(exe)

	storePath, err := store.Path()
	if err != nil {
		return "", err
	}

	binDir, err := store.Bin(storePath)
	if err != nil {
		return "", err
	}

	// NOTE: The expected directory is resolved too; '$GDENV_HOME' may itself be
	// reached through a symlink (e.g. a symlinked '$HOME'), in which case a
	// naive comparison would reject a perfectly valid installation.
	want := resolveSymlinks(binDir)

	if !samePath(filepath.Dir(exe), want) {
		return "", fmt.Errorf(
			"%w: running '%s', expected it within '%s'",
			ErrUnmanaged, exe, binDir,
		)
	}

	return binDir, nil
}

/* ------------------------ Function: resolveSymlinks ----------------------- */

// resolveSymlinks expands any symbolic links within a path, returning it
// unchanged if that cannot be done.
//
// NOTE: Resolution is best-effort because the path may not exist yet; in that
// case the unresolved form is the best available answer.
func resolveSymlinks(path string) string {
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return path
	}

	return resolved
}

/* ----------------------------- Function: samePath ------------------------- */

// samePath reports whether two paths refer to the same location.
func samePath(a, b string) bool {
	a, b = filepath.Clean(a), filepath.Clean(b)

	// NOTE: Windows filepaths are case-insensitive.
	if runtime.GOOS == osWindows {
		return strings.EqualFold(a, b)
	}

	return a == b
}

/* -------------------------------------------------------------------------- */
/*                               Function: Apply                              */
/* -------------------------------------------------------------------------- */

// Apply downloads the specified version and replaces the 'gdenv' and 'godot'
// binaries within 'binDir'.
func Apply(ctx context.Context, binDir, v string) error {
	// NOTE: The version is normalized here as well as by the command, so that
	// the release asset is never requested under a name like '0.7.0' which would
	// surface as an unexplained 404.
	v = Normalize(v)
	if !IsValidVersion(v) {
		return fmt.Errorf("%w: %s", ErrInvalidVersion, v)
	}

	target, err := DetectTarget()
	if err != nil {
		return err
	}

	if err := verifyWritable(binDir); err != nil {
		return err
	}

	unlock, err := acquireLock(filepath.Dir(binDir))
	if err != nil {
		return err
	}

	defer unlock()

	// Sweep artifacts left behind by a previous update before adding more.
	//
	// NOTE: This must happen *after* the lock is acquired; otherwise a second
	// process could remove the staging directory of an update already running.
	cleanStale(binDir)

	// NOTE: The staging directory is created *within* the install directory so
	// that the final renames stay on one filesystem. This deviates from the
	// 'os.MkdirTemp("", ...)' used elsewhere, where atomicity does not matter.
	tmp, err := os.MkdirTemp(binDir, prefixTempDir+"*")
	if err != nil {
		return err
	}

	defer os.RemoveAll(tmp)

	pathArchive, err := fetchArchive(ctx, target, v, tmp)
	if err != nil {
		return err
	}

	dirExtract := filepath.Join(tmp, "extract")
	if err := os.Mkdir(dirExtract, osutil.ModeUserRWX); err != nil {
		return err
	}

	if err := extractArchive(ctx, target, pathArchive, dirExtract); err != nil {
		return err
	}

	names := target.Binaries()

	found, err := locateBinaries(dirExtract, names)
	if err != nil {
		return err
	}

	return replaceAll(binDir, names, found)
}

/* -------------------------- Function: verifyWritable ---------------------- */

// verifyWritable confirms that files can be created within the directory, so
// that a read-only installation fails before anything is downloaded.
func verifyWritable(dir string) error {
	f, err := os.CreateTemp(dir, prefixTempDir+"check-*")
	if err != nil {
		return fmt.Errorf("%w: %s: %w", ErrNotWritable, dir, err)
	}

	// NOTE: 'cleanStale' does not sweep this file - it removes superseded
	// binaries and staging *directories* - so it must be removed on every path
	// out of this function.
	defer os.Remove(f.Name())

	return f.Close()
}

/* -------------------------- Function: acquireLock ------------------------- */

// acquireLock takes an advisory lock preventing concurrent updates. The
// returned function releases it.
//
// NOTE: Reclaiming an abandoned lock is not atomic - between the staleness
// check and the removal below, another process may already have replaced it,
// in which case both would proceed. Doing this properly needs more than the
// filesystem primitives used here, and the window requires two updates racing
// within moments of each other against a lock abandoned for 'lockStaleAfter',
// so it is accepted rather than engineered around.
func acquireLock(dir string) (func(), error) {
	path := filepath.Join(dir, lockName)

	unlock, err := createLock(path)
	if err == nil {
		return unlock, nil
	}

	if !errors.Is(err, fs.ErrExist) {
		return nil, err
	}

	// An abandoned lock - left by an interrupted or killed update - should not
	// block future updates forever.
	info, err := os.Stat(path)
	if err != nil {
		// NOTE: The lock was released between the two calls above, so there is
		// nothing to reclaim; simply try again.
		if errors.Is(err, fs.ErrNotExist) {
			return createLock(path)
		}

		return nil, err
	}

	if time.Since(info.ModTime()) < lockStaleAfter {
		return nil, ErrLockHeld
	}

	log.Debugf("removing stale update lock: %s", path)

	if err := os.Remove(path); err != nil {
		return nil, err
	}

	return createLock(path)
}

/* --------------------------- Function: createLock ------------------------- */

func createLock(path string) (func(), error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, osutil.ModeUserRW)
	if err != nil {
		return nil, err
	}

	if err := f.Close(); err != nil {
		return nil, err
	}

	return func() {
		if err := os.Remove(path); err != nil {
			log.Debugf("could not remove update lock: %v", err)
		}
	}, nil
}
