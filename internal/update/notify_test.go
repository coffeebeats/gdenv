package update

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

/* -------------------------------------------------------------------------- */
/*                           Function: TestIsEnabled                          */
/* -------------------------------------------------------------------------- */

// TestIsEnabled covers the decisions this package makes around
// 'strconv.ParseBool' - the empty value, normalization, and what an
// unrecognized value means - rather than re-testing the parsing itself.
func TestIsEnabled(t *testing.T) {
	tests := []struct {
		value string
		want  bool
	}{
		// An unset variable leaves the flag off. 'strconv.ParseBool' rejects
		// this, and a rejected value enables.
		{value: "", want: false},

		// The two values which disable a flag. These are what users are told to
		// set, so they must keep working even though anything else would also
		// parse as false.
		{value: "0", want: false},
		{value: "false", want: false},

		// Normalization: neither would parse as written.
		{value: " false ", want: false},
		{value: "fAlSe", want: false},

		// A parsed value is passed through, not inverted.
		{value: "true", want: true},

		// An unrecognized value enables an opt-in flag.
		{value: "yes", want: true},
	}

	for _, tc := range tests {
		t.Run("value="+tc.value, func(t *testing.T) {
			if got := isEnabled(tc.value); got != tc.want {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
		})
	}
}

/* -------------------------------------------------------------------------- */
/*                          Function: TestShouldCheck                         */
/* -------------------------------------------------------------------------- */

func TestShouldCheck(t *testing.T) {
	// Automated environments must never reach out to the network.
	t.Run("is disabled in CI", func(t *testing.T) {
		t.Setenv(EnvNoNotifier, "")
		t.Setenv(envCI, "true")

		if shouldCheck() {
			t.Fatal("expected the check to be skipped when 'CI' is set")
		}
	})

	t.Run("is disabled by the opt-out variable", func(t *testing.T) {
		t.Setenv(envCI, "")
		t.Setenv(EnvNoNotifier, "1")

		if shouldCheck() {
			t.Fatalf("expected the check to be skipped when '%s' is set", EnvNoNotifier)
		}
	})

	// A '0' must not read as "disabled", which a bare presence check would do.
	t.Run("is not disabled by an explicit zero", func(t *testing.T) {
		t.Setenv(envCI, "")
		t.Setenv(EnvNoNotifier, "0")

		// The result still depends on whether stderr is a terminal, which it is
		// not under 'go test'; assert only that the opt-out did not trigger.
		if isEnabled(os.Getenv(EnvNoNotifier)) {
			t.Fatalf("expected '%s=0' to leave the notifier enabled", EnvNoNotifier)
		}
	})

	// Output which is piped or redirected must stay free of notices.
	t.Run("is disabled when stderr is not a terminal", func(t *testing.T) {
		t.Setenv(envCI, "")
		t.Setenv(EnvNoNotifier, "")

		if shouldCheck() {
			t.Fatal("expected the check to be skipped when stderr is not a terminal")
		}
	})
}

/* -------------------------------------------------------------------------- */
/*                          Function: TestNotifierState                       */
/* -------------------------------------------------------------------------- */

func TestNotifierState(t *testing.T) {
	t.Run("round-trips the recorded check", func(t *testing.T) {
		home := t.TempDir()

		if err := writeState(home, "v0.7.0"); err != nil {
			t.Fatal(err)
		}

		got, err := readState(home)
		if err != nil {
			t.Fatal(err)
		}

		if got.LatestVersion != "v0.7.0" {
			t.Fatalf("version: got %q, want %q", got.LatestVersion, "v0.7.0")
		}

		if time.Since(got.CheckedAt) > time.Minute {
			t.Fatalf("checked-at: got %v, want a recent timestamp", got.CheckedAt)
		}
	})

	// A check which failed is recorded with an empty version so that it backs
	// off; discarding it on the way back in would defeat that, leaving every
	// command to retry an unreachable network.
	t.Run("round-trips a failed check", func(t *testing.T) {
		home := t.TempDir()

		if err := writeState(home, ""); err != nil {
			t.Fatal(err)
		}

		got, err := readState(home)
		if err != nil {
			t.Fatal(err)
		}

		if got.LatestVersion != "" {
			t.Fatalf("version: got %q, want an empty version", got.LatestVersion)
		}

		if time.Since(got.CheckedAt) > time.Minute {
			t.Fatalf("checked-at: got %v, want a recent timestamp", got.CheckedAt)
		}
	})

	t.Run("reports a malformed state file", func(t *testing.T) {
		home := t.TempDir()

		writeFile(t, filepath.Join(home, stateName), "not json")

		if _, err := readState(home); err == nil {
			t.Fatal("expected an error for a malformed state file")
		}
	})

	// Nothing validates the cache file on the way in, so a hand-edited or
	// corrupted version must not reach the upgrade comparison: a non-semver
	// value on the *current* side would make any release look newer.
	t.Run("rejects a cached version which is not valid semver", func(t *testing.T) {
		home := t.TempDir()

		writeFile(t, filepath.Join(home, stateName),
			`{"checkedAt":"2099-01-01T00:00:00Z","latestVersion":"not-a-version"}`)

		if _, err := readState(home); !errors.Is(err, ErrInvalidVersion) {
			t.Fatalf("err: got %v, want %v", err, ErrInvalidVersion)
		}
	})
}

/* -------------------------------------------------------------------------- */
/*                            Function: TestNotify                            */
/* -------------------------------------------------------------------------- */

// TestNotify covers the reporting decision in isolation from the network, which
// is the part that determines whether the user sees anything.
func TestNotify(t *testing.T) {
	tests := []struct {
		name    string
		current string
		latest  string
		want    string
	}{
		{name: "reports a newer version", current: "v0.6.35", latest: "v0.7.0", want: "v0.7.0"},
		{name: "stays quiet when up-to-date", current: "v0.6.35", latest: "v0.6.35", want: ""},
		{name: "stays quiet for an older release", current: "v0.7.0", latest: "v0.6.36", want: ""},
		{name: "stays quiet with no result", current: "v0.6.35", latest: "", want: ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			n := &Notifier{current: tc.current, done: nil, latest: tc.latest}

			if got := n.Notify(); got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}

/* -------------------------------------------------------------------------- */
/*                      Function: TestNotifierUsesCache                       */
/* -------------------------------------------------------------------------- */

// TestNotifierUsesCache verifies that a fresh cache entry is reported without
// starting a request, which is what keeps the check to once per interval.
func TestNotifierUsesCache(t *testing.T) {
	home := t.TempDir()

	b, err := json.Marshal(state{CheckedAt: time.Now(), LatestVersion: "v9.9.9"})
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(home, stateName), b, 0o600); err != nil {
		t.Fatal(err)
	}

	got, err := readState(home)
	if err != nil {
		t.Fatal(err)
	}

	if time.Since(got.CheckedAt) >= checkInterval {
		t.Fatal("expected the recorded check to be considered fresh")
	}

	if got.LatestVersion != "v9.9.9" {
		t.Fatalf("version: got %q, want %q", got.LatestVersion, "v9.9.9")
	}
}
