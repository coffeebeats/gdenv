package update

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
)

const (
	extensionStale = ".old"
	prefixTempDir  = ".update-"
)

// rename is the rename implementation used when replacing binaries. It exists
// as a variable so that tests can simulate a rename failure on any platform.
//
// NOTE: Tests reassign this, so they must not call 't.Parallel'.
//
//nolint:gochecknoglobals
var rename = os.Rename

/* --------------------------- Function: replaceAll ------------------------- */

// replaceAll replaces each named binary within 'binDir', rolling back any which
// already succeeded should a later one fail.
func replaceAll(binDir string, names []string, found map[string]string) error {
	restores := make([]func() error, 0, len(names))

	for _, name := range names {
		// NOTE: 'locateBinaries' guarantees an entry for every name; guard anyway,
		// so a future caller gets this rather than a rename of the empty path.
		src, ok := found[name]
		if !ok {
			return fmt.Errorf("%w: %s", ErrMissingBinary, name)
		}

		restore, err := replaceBinary(src, filepath.Join(binDir, name))
		if err != nil {
			rollback(restores)

			return fmt.Errorf("failed to replace '%s': %w", name, err)
		}

		log.Debugf("replaced binary: %s", name)

		restores = append(restores, restore)
	}

	commitReplacements(binDir, names)

	return nil
}

/* --------------------------- Function: stalePath -------------------------- */

// stalePath returns the path at which a binary's superseded copy is kept.
func stalePath(target string) string {
	return target + extensionStale
}

/* -------------------------------------------------------------------------- */
/*                          Function: replaceBinary                           */
/* -------------------------------------------------------------------------- */

// replaceBinary replaces the file at 'target' with the file at 'src' and
// returns a function which restores the original.
//
// The original is moved aside to a sibling '.old' file before the replacement
// is moved into place. Doing so uniformly - rather than overwriting 'target'
// directly - serves two purposes:
//
//  1. Windows refuses to overwrite the executable of a running process, but
//     *does* permit renaming it. This is what allows 'gdenv' to replace itself.
//  2. It keeps the original recoverable, so that a failure partway through
//     replacing several binaries can be rolled back.
//
// The cost is a brief window during which 'target' does not exist. That is
// preferable to the alternative, where the fast path overwrites the original
// and leaves nothing to roll back to.
func replaceBinary(src, target string) (func() error, error) {
	stale := stalePath(target)

	hasTarget, err := moveAside(target, stale)
	if err != nil {
		return nil, err
	}

	if err := rename(src, target); err != nil {
		if !hasTarget {
			return nil, err
		}

		// Put the original back before reporting the failure.
		if errRestore := rename(stale, target); errRestore != nil {
			return nil, errors.Join(err, errRestore)
		}

		return nil, err
	}

	// NOTE: There was nothing at 'target' beforehand, so undoing the
	// replacement means removing what was put there, not restoring anything.
	if !hasTarget {
		return func() error { return removeIfPresent(target) }, nil
	}

	return func() error { return restoreOriginal(target, stale) }, nil
}

/* ---------------------------- Function: moveAside ------------------------- */

// moveAside clears any leftover stale file and moves an existing target out of
// the way, reporting whether there was in fact a target to move.
//
// NOTE: A missing target is not an error; the binary may have been removed, in
// which case there is simply nothing to restore.
func moveAside(target, stale string) (bool, error) {
	// Remove any leftover from a previous update, so that moving the current
	// binary aside cannot fail because the destination already exists.
	if err := os.Remove(stale); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return false, err
	}

	if _, err := os.Stat(target); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil
		}

		return false, err
	}

	if err := rename(target, stale); err != nil {
		return false, err
	}

	return true, nil
}

/* ------------------------- Function: restoreOriginal ---------------------- */

// restoreOriginal undoes a replacement performed by 'replaceBinary', moving the
// superseded binary back into place.
func restoreOriginal(target, stale string) error {
	if err := removeIfPresent(target); err != nil {
		return err
	}

	return rename(stale, target)
}

/* ------------------------- Function: removeIfPresent ---------------------- */

// removeIfPresent deletes a file, treating an already-absent one as success.
func removeIfPresent(path string) error {
	if err := os.Remove(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	return nil
}

/* ----------------------- Function: commitReplacements --------------------- */

// commitReplacements removes the superseded binaries left behind by a
// successful replacement.
//
// NOTE: This is best-effort; on Windows the '.old' file belonging to the
// currently running executable is still locked, so it is swept by a later
// update instead.
func commitReplacements(binDir string, names []string) {
	for _, name := range names {
		stale := stalePath(filepath.Join(binDir, name))

		if err := os.Remove(stale); err != nil && !errors.Is(err, fs.ErrNotExist) {
			log.Debugf("could not remove superseded binary: %s: %v", name, err)
		}
	}
}

/* ---------------------------- Function: rollback -------------------------- */

// rollback restores previously replaced binaries, reporting - but not returning
// - any failure to do so; the error which triggered the rollback is the one
// worth surfacing.
func rollback(restores []func() error) {
	for _, restore := range restores {
		if err := restore(); err != nil {
			log.Errorf("failed to restore the previous binary: %v", err)
			log.Error("the installation may be broken; reinstall by following " +
				"https://github.com/coffeebeats/gdenv/blob/main/docs/installation.md")
		}
	}
}

/* -------------------------------------------------------------------------- */
/*                            Function: cleanStale                            */
/* -------------------------------------------------------------------------- */

// cleanStale removes artifacts left behind by a previous update.
//
// NOTE: This is best-effort. The '.old' file belonging to the currently running
// binary cannot be removed on Windows, so it is swept by a later invocation.
func cleanStale(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Debugf("could not scan for stale update artifacts: %s: %v", dir, err)

		return
	}

	for _, entry := range entries {
		name := entry.Name()

		isStaleBinary := !entry.IsDir() && strings.HasSuffix(name, extensionStale)
		isStaleTempDir := entry.IsDir() && strings.HasPrefix(name, prefixTempDir)

		if !isStaleBinary && !isStaleTempDir {
			continue
		}

		if err := os.RemoveAll(filepath.Join(dir, name)); err != nil {
			log.Debugf("could not remove stale update artifact: %s: %v", name, err)
		}
	}
}
