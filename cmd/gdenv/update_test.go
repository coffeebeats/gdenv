package main

import (
	"errors"
	"strings"
	"testing"

	"github.com/urfave/cli/v2"
)

/* -------------------------------------------------------------------------- */
/*                          Function: TestUpdateUsage                         */
/* -------------------------------------------------------------------------- */

// TestUpdateUsage covers the flag and argument combinations the command
// rejects. Each is caught before a request is made, so none reach the network.
func TestUpdateUsage(t *testing.T) {
	const current = "v0.6.35"

	tests := []struct {
		name string
		args []string
		want error
	}{
		{
			name: "rejects '--check' with '--force'",
			args: []string{"gdenv", "update", "--check", "--force"},
			want: ErrUpdateUsageCheckAndForce,
		},
		{
			name: "rejects '--check' with a version",
			args: []string{"gdenv", "update", "--check", "v0.7.0"},
			want: ErrUpdateUsageCheckAndVersion,
		},
		{
			name: "rejects a version which is not semver",
			args: []string{"gdenv", "update", "not-a-version"},
			want: ErrUpdateUsageInvalidVersion,
		},

		// A named version is installed without consulting the latest release,
		// so these two are the only guard against an unintended downgrade or
		// reinstall.
		{
			name: "rejects an older version without '--force'",
			args: []string{"gdenv", "update", "v0.5.0"},
			want: ErrUpdateUsageNotNewer,
		},
		{
			name: "rejects the running version without '--force'",
			args: []string{"gdenv", "update", current},
			want: ErrUpdateUsageNotNewer,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			app := &cli.App{
				Name:     "gdenv",
				Version:  current,
				Commands: []*cli.Command{NewUpdate()},
			}

			err := app.Run(tc.args)
			if !errors.Is(err, tc.want) {
				t.Fatalf("err: got %v, want %v", err, tc.want)
			}

			// The error must reach 'main' as a 'UsageError' so that the
			// command's usage is printed alongside it.
			var usageErr UsageError
			if !errors.As(err, &usageErr) {
				t.Fatalf("err: got %T, want a UsageError", err)
			}
		})
	}
}

/* -------------------------------------------------------------------------- */
/*                         Function: TestUpdateNotice                         */
/* -------------------------------------------------------------------------- */

func TestUpdateNotice(t *testing.T) {
	lines := updateNotice("v0.6.0", "v0.6.35")

	const wantLines = 3
	if len(lines) != wantLines {
		t.Fatalf("lines: got %d, want %d", len(lines), wantLines)
	}

	// Both versions must appear, so the notice shows the transition rather than
	// only the version being offered.
	for _, want := range []string{noticeIcon, "v0.6.0", "v0.6.35"} {
		if !strings.Contains(lines[0], want) {
			t.Errorf("header %q does not contain %q", lines[0], want)
		}
	}

	if !strings.Contains(lines[1], "gdenv update") {
		t.Errorf("hint %q does not name the command", lines[1])
	}

	if !strings.Contains(lines[2], "releases/tag/v0.6.35") {
		t.Errorf("link %q does not point at the release notes", lines[2])
	}
}
