package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v2"

	"github.com/coffeebeats/gdenv/internal/update"
)

const (
	cmdUpdate = "update"

	flagCheck = "check"

	// noticeIcon prefixes the "update available" notice so that it stands out
	// from the surrounding command output.
	noticeIcon = "✨"
)

var (
	ErrUpdateUsageCheckAndForce   = errors.New("cannot specify both '-c/--check' and '-f/--force'")
	ErrUpdateUsageCheckAndVersion = errors.New("cannot specify 'VERSION' with '-c/--check'")
	ErrUpdateUsageInvalidVersion  = errors.New("invalid version specified")
	ErrUpdateUsageNotNewer        = errors.New("version is not newer; specify '-f/--force' to install it anyway")
)

/* --------------------------- Function: NewUpdate -------------------------- */

// A 'urfave/cli' command to update the 'gdenv' installation itself.
func NewUpdate() *cli.Command {
	return &cli.Command{
		Name:     cmdUpdate,
		Category: categoryUtilities,

		Usage:     "update gdenv itself to the latest released version",
		UsageText: "gdenv update [OPTIONS] [VERSION]",

		Flags: []cli.Flag{
			newVerboseFlag(),

			&cli.BoolFlag{
				Name:    flagCheck,
				Aliases: []string{"c"},
				Usage:   "check for a new version without installing it",
			},
			&cli.BoolFlag{
				Name:    flagForce,
				Aliases: []string{"f"},
				Usage:   "install even if the version is not newer",
			},
		},

		Action: func(c *cli.Context) error {
			target, err := validateUpdateArgs(c)
			if err != nil {
				return err
			}

			if c.Bool(flagCheck) {
				return checkForUpdate(c.Context, c.App.Version)
			}

			return applyUpdate(c.Context, c.App.Version, target, c.Bool(flagForce))
		},
	}
}

/* ---------------------- Function: validateUpdateArgs ---------------------- */

// validateUpdateArgs checks the command's flags and argument against one
// another, returning the version to install; an empty string means the latest.
func validateUpdateArgs(c *cli.Context) (string, error) {
	if c.Bool(flagCheck) && c.Bool(flagForce) {
		return "", UsageError{ctx: c, err: ErrUpdateUsageCheckAndForce}
	}

	target := update.Normalize(c.Args().First())
	if target == "" {
		return "", nil
	}

	if c.Bool(flagCheck) {
		return "", UsageError{ctx: c, err: ErrUpdateUsageCheckAndVersion}
	}

	if !update.IsValidVersion(target) {
		return "", UsageError{
			ctx: c,
			err: fmt.Errorf("%w: %s", ErrUpdateUsageInvalidVersion, c.Args().First()),
		}
	}

	// NOTE: A named 'VERSION' is installed without consulting the latest
	// release, so this is the only place a downgrade - or a reinstall of the
	// running version - can be caught. Requiring '-f/--force' for it matches
	// what updating to the latest version already requires.
	if !c.Bool(flagForce) && !update.IsUpgrade(c.App.Version, target) {
		return "", UsageError{
			ctx: c,
			err: fmt.Errorf("%w: %s", ErrUpdateUsageNotNewer, target),
		}
	}

	return target, nil
}

/* ------------------------ Function: checkForUpdate ------------------------ */

// checkForUpdate reports the latest available version without installing it.
//
// NOTE: This always exits successfully; an available update is information, not
// an error.
func checkForUpdate(ctx context.Context, current string) error {
	status, err := update.Check(ctx, current)
	if err != nil {
		return err
	}

	if !status.IsUpgrade() {
		reportUpToDate(status.Current)

		return nil
	}

	log.Printf("gdenv %s is available (currently %s)", status.Latest, status.Current)
	log.Print("run 'gdenv update' to install it")

	return nil
}

/* -------------------------- Function: applyUpdate ------------------------- */

// applyUpdate downloads and installs the specified version, defaulting to the
// latest published one.
func applyUpdate(ctx context.Context, current, target string, force bool) error {
	if target == "" {
		status, err := update.Check(ctx, current)
		if err != nil {
			return err
		}

		if !status.IsUpgrade() && !force {
			reportUpToDate(status.Current)

			return nil
		}

		target = status.Latest
	}

	// NOTE: The install directory is resolved only once there is something to
	// install, so that an installation gdenv does not manage is not reported as
	// unmanaged when it was already up-to-date.
	binDir, err := update.ManagedBinDir()
	if err != nil {
		return err
	}

	if err := update.Apply(ctx, binDir, target); err != nil {
		return err
	}

	log.Printf("updated gdenv to %s", target)

	return nil
}

/* ------------------------ Function: reportUpToDate ------------------------ */

// reportUpToDate reports that the running version is the latest one.
func reportUpToDate(current string) {
	log.Printf("gdenv %s is up-to-date", current)
}

/* -------------------------------------------------------------------------- */
/*                           Function: reportUpdate                           */
/* -------------------------------------------------------------------------- */

// reportUpdate surfaces a newly released version once the command has finished.
// When automatic updates are opted into, the new version is installed instead
// of being announced.
//
// NOTE: This must never alter the command's output contract or exit code; all
// failures are reported as warnings and otherwise ignored.
func reportUpdate(notifier *update.Notifier, interrupted bool, exitCode int) {
	// Exit promptly when interrupted rather than waiting on the check.
	if notifier == nil || interrupted {
		return
	}

	latest := notifier.Notify()
	if latest == "" {
		return
	}

	if !update.IsAutoUpdate() {
		// Separate the notice from the command's own output.
		log.Print("")

		for _, line := range updateNotice(notifier.Current(), latest) {
			log.Print(line)
		}

		return
	}

	// Don't start a download on the way out of a failed command.
	if exitCode != 0 {
		return
	}

	autoUpdate(latest)
}

/* -------------------------- Function: autoUpdate -------------------------- */

// autoUpdate installs a newly released version without being asked to.
//
// NOTE: The command has already finished by this point, so nothing here may
// change its outcome; every failure is reported and then dropped.
func autoUpdate(latest string) {
	binDir, err := update.ManagedBinDir()
	if err != nil {
		log.Debugf("skipping automatic update: %v", err)

		return
	}

	log.Infof("automatically updating gdenv to %s", latest)

	// NOTE: The command's context has already been cancelled, and the download
	// is deliberately left unbounded; when to stop waiting on it is the user's
	// call rather than a timeout's.
	if err := update.Apply(context.Background(), binDir, latest); err != nil {
		if errors.Is(err, update.ErrLockHeld) {
			log.Debug("skipping automatic update; another update is in progress")

			return
		}

		log.Warnf("automatic update failed: %v", err)

		return
	}

	log.Infof("updated gdenv to %s", latest)
}

/* -------------------------------------------------------------------------- */
/*                           Function: updateNotice                           */
/* -------------------------------------------------------------------------- */

// updateNotice renders the "new version available" message as a set of lines.
//
// NOTE: This is only ever shown on an interactive terminal (see the notifier's
// own gating), so styling it is safe.
func updateNotice(current, latest string) []string {
	var (
		styleHeader = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.ANSIColor(colorYellowBright))
		styleOld    = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(colorRedBright))
		styleNew    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.ANSIColor(colorGreenBright))
		styleCmd    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.ANSIColor(colorWhiteBright))
		styleLink   = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(colorCyanBright))
		styleFaint  = lipgloss.NewStyle().Faint(true)
	)

	// Continuation lines are indented to align beneath the header's text.
	const indent = "   "

	return []string{
		fmt.Sprintf(
			"%s %s %s %s %s",
			noticeIcon,
			styleHeader.Render("A new release of gdenv is available:"),
			styleOld.Render(current),
			styleFaint.Render("→"),
			styleNew.Render(latest),
		),
		fmt.Sprintf("%sTo upgrade, run: %s", indent, styleCmd.Render("gdenv update")),
		indent + styleLink.Render(update.ReleaseNotesURL(latest)),
	}
}
