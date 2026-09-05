package update

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/mattn/go-isatty"

	"github.com/coffeebeats/gdenv/internal/osutil"
	"github.com/coffeebeats/gdenv/pkg/store"
)

const (
	// EnvAutoUpdate opts in to applying updates automatically.
	//
	// NOTE: This builds on the notifier's result, so 'EnvNoNotifier' disables it
	// as well; there is no separate check for an update to be applied from.
	EnvAutoUpdate = "GDENV_AUTO_UPDATE"
	// EnvNoNotifier disables the periodic check for new versions.
	EnvNoNotifier = "GDENV_NO_UPDATE_NOTIFIER"

	envCI = "CI"

	stateName = ".update-check.json"

	// checkInterval is how long a resolved version is cached before the check
	// is repeated.
	checkInterval = 24 * time.Hour

	// checkTimeout bounds the background check. 'internal/client' sets no
	// timeout of its own, so without this an unreachable network could hold up
	// the command's exit.
	checkTimeout = 2 * time.Second

	// waitTimeout is how long the command will wait at exit for an in-flight
	// check. Exceeding it costs only the notice, since the check persists its
	// own result for the next invocation.
	waitTimeout = time.Second
)

/* -------------------------------------------------------------------------- */
/*                                Struct: state                               */
/* -------------------------------------------------------------------------- */

// state is the on-disk record of the most recent version check.
//
// NOTE: 'CheckedAt' records the most recent *attempt*, not the most recent
// success. A failed check is recorded with an empty 'LatestVersion' so that an
// unreachable network backs off, rather than being retried by every command.
type state struct {
	CheckedAt     time.Time `json:"checkedAt"`
	LatestVersion string    `json:"latestVersion"`
}

/* ---------------------------- Function: statePath ------------------------- */

// statePath returns the path at which the version check's result is cached.
func statePath(storePath string) string {
	return filepath.Join(storePath, stateName)
}

/* ---------------------------- Function: readState ------------------------- */

func readState(storePath string) (state, error) {
	var s state

	b, err := os.ReadFile(statePath(storePath))
	if err != nil {
		return s, err
	}

	if err := json.Unmarshal(b, &s); err != nil {
		return s, err
	}

	// NOTE: Nothing validated this file on the way in - it may have been edited
	// or corrupted - so a version which is not valid semver is discarded rather
	// than compared against the running one. An empty version is expected; it
	// records a check which was attempted but did not succeed.
	if s.LatestVersion != "" && !IsValidVersion(s.LatestVersion) {
		return state{}, fmt.Errorf("%w: %s", ErrInvalidVersion, s.LatestVersion)
	}

	return s, nil
}

/* --------------------------- Function: writeState ------------------------- */

// NOTE: The write is not atomic. A torn file fails to parse in 'readState' and
// is discarded, which triggers a fresh check and overwrites it - so the cache
// heals itself without needing a temporary file and a rename.
func writeState(storePath, latest string) error {
	b, err := json.Marshal(state{CheckedAt: time.Now(), LatestVersion: latest})
	if err != nil {
		return err
	}

	if err := os.WriteFile(statePath(storePath), b, osutil.ModeUserRW); err != nil {
		// The store may not exist yet; that is not worth reporting.
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}

		return err
	}

	return nil
}

/* -------------------------------------------------------------------------- */
/*                               Struct: Notifier                             */
/* -------------------------------------------------------------------------- */

// A Notifier performs a cached, best-effort check for newer versions of
// 'gdenv' alongside the command being run.
type Notifier struct {
	current string
	done    chan struct{}
	latest  string
}

/* --------------------------- Function: NewNotifier ------------------------ */

// NewNotifier starts a background version check for the running command, unless
// checking is disabled. The returned 'Notifier' is always safe to use; a
// disabled one simply reports nothing.
//
// NOTE: Nothing here may affect the command's output or exit code. Every error
// is deliberately swallowed.
func NewNotifier(ctx context.Context, current string) *Notifier {
	n := &Notifier{current: Normalize(current), done: nil, latest: ""}

	if !shouldCheck() {
		return n
	}

	storePath, err := store.Path()
	if err != nil {
		return n
	}

	// A cached result is used immediately; only a stale cache triggers a
	// request.
	if s, err := readState(storePath); err == nil && time.Since(s.CheckedAt) < checkInterval {
		n.latest = s.LatestVersion

		return n
	}

	n.done = make(chan struct{})

	// NOTE: Cancellation is dropped deliberately. The caller's context is
	// cancelled as soon as the command finishes, which would abort this check
	// before it ever completed; 'checkTimeout' bounds it instead.
	ctxCheck := context.WithoutCancel(ctx)

	go func() {
		defer close(n.done)

		ctx, cancel := context.WithTimeout(ctxCheck, checkTimeout)
		defer cancel()

		latest, err := LatestVersion(ctx)
		if err != nil {
			log.Debugf("update check failed: %v", err)

			latest = ""
		}

		n.latest = latest

		// Record the attempt even if the command exits before reading it, and
		// even if it failed: the next invocation can then report without a
		// request, and an unreachable network is not retried by every command.
		if err := writeState(storePath, latest); err != nil {
			log.Debugf("could not record update check: %v", err)
		}
	}()

	return n
}

/* ------------------------------ Method: Notify ---------------------------- */

// Notify waits briefly for any in-flight check and then reports whether a newer
// version is available. It returns the newer version, or an empty string.
//
// NOTE: This must be called from the same deferred block which exits the
// process, so that it still runs when a command fails.
func (n *Notifier) Notify() string {
	if n.done != nil {
		select {
		case <-n.done:
		case <-time.After(waitTimeout):
			// NOTE: This branch must not read 'n.latest'. The check may still be
			// writing it; receiving above is what orders that write before the
			// read below.
			return ""
		}
	}

	if !IsUpgrade(n.current, n.latest) {
		return ""
	}

	return n.latest
}

/* ------------------------------ Method: Current --------------------------- */

// Current returns the normalized version of the running binary.
func (n *Notifier) Current() string {
	return n.current
}

/* ------------------------- Function: IsAutoUpdate ------------------------- */

// IsAutoUpdate reports whether updates should be applied without being asked
// for.
func IsAutoUpdate() bool {
	return isEnabled(os.Getenv(EnvAutoUpdate))
}

/* --------------------------- Function: shouldCheck ------------------------ */

// shouldCheck reports whether a version check is appropriate for this run.
func shouldCheck() bool {
	if isEnabled(os.Getenv(EnvNoNotifier)) {
		return false
	}

	// Never reach out to the network from an automated environment.
	//
	// NOTE: Unlike the variables above, any value counts here - including
	// 'false', which some tools set explicitly. Skipping the check in a job
	// which is not really CI costs nothing; running it there could.
	if os.Getenv(envCI) != "" {
		return false
	}

	// Only interactive use gets a notice; piped output should stay clean.
	return isatty.IsTerminal(os.Stderr.Fd()) || isatty.IsCygwinTerminal(os.Stderr.Fd())
}

/* --------------------------- Function: isEnabled -------------------------- */

// isEnabled interprets an environment variable as a boolean flag. Any value
// other than the empty string, '0', or 'false' enables it.
//
// NOTE: These are opt-in flags, so an unrecognized value - 'yes', 'on' - enables
// rather than disables; only an explicit false switches one off. That inverts
// how 'strconv.ParseBool' is used elsewhere (see 'version.LabelDefault'), where
// an unparseable value selects the default instead.
func isEnabled(value string) bool {
	// Normalizing first widens what 'strconv.ParseBool' accepts to any casing,
	// and keeps a stray space from reading as an unrecognized - and therefore
	// enabling - value.
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return false
	}

	enabled, err := strconv.ParseBool(value)
	if err != nil {
		return true
	}

	return enabled
}
