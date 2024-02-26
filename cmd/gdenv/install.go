package main

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v2"

	"github.com/coffeebeats/gdenv/pkg/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/pkg/godot/platform"
	"github.com/coffeebeats/gdenv/pkg/godot/version"
	"github.com/coffeebeats/gdenv/pkg/install"
	"github.com/coffeebeats/gdenv/pkg/pin"
	"github.com/coffeebeats/gdenv/pkg/store"
)

var (
	ErrInstallUsageGlobalAndPath   = errors.New("cannot specify both '-g/--global' and '-p/--path'")
	ErrInstallUsageGlobalAndSource = errors.New("cannot specify both '-g/--global' and '-s/--source'")
)

// A 'urfave/cli' command to download and cache a specific version of Godot.
func NewInstall() *cli.Command { //nolint:funlen
	return &cli.Command{
		Name:     "install",
		Category: "Install",

		Aliases: []string{"i"},

		Usage: "download and cache a specific version of Godot; " +
			"if 'VERSION' is omitted then the version is resolved using '-g', '-p', or '$PWD'",
		UsageText: "gdenv install [OPTIONS] [VERSION]",

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "forcibly overwrite an existing cache entry",
			},
			&cli.BoolFlag{
				Name:    "global",
				Aliases: []string{"g"},
				Usage:   "update the global pin (if 'VERSION' is specified) or resolve 'VERSION' from the global pin",
			},
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Usage:   "resolve the pinned 'VERSION' at 'PATH'",
			},
			&cli.BoolFlag{
				Name:    "source",
				Aliases: []string{"s", "src"},
				Usage:   "install source code instead of an executable (cannot be used with '-g')",
			},
		},

		Action: func(c *cli.Context) error {
			// Validate flag options.
			if c.IsSet("global") && c.IsSet("path") {
				return UsageError{ctx: c, err: ErrPinUsageGlobalAndPath}
			}
			if c.IsSet("global") && c.IsSet("source") {
				return UsageError{ctx: c, err: ErrInstallUsageGlobalAndSource}
			}

			v, err := resolveVersionFromInput(c)
			if err != nil {
				return err
			}

			storePath, err := touchStore()
			if err != nil {
				return err
			}

			log.Debugf("using store at path: %s", storePath)

			if c.Bool("source") {
				return install.Source(c.Context, storePath, v, c.Bool("force"))
			}

			if err := installExecutable(c.Context, storePath, v, c.Bool("force")); err != nil {
				return err
			}

			if !c.Bool("global") {
				return nil
			}

			return writePin(storePath, storePath, v)
		},
	}
}

/* ----------------------- Function: installExecutable ---------------------- */

// Installs the specified executable version to the store, but only if needed.
func installExecutable(
	ctx context.Context,
	storePath string,
	v version.Version,
	force bool,
) error {
	// Define the host 'Platform'.
	p, err := platform.Detect()
	if err != nil {
		return err
	}

	// Define the target 'Executable'.
	ex := executable.New(v, p)

	ok, err := store.Has(storePath, ex)
	if err != nil {
		return err
	}

	if ok && !force {
		log.Info("skipping installation; version already found")

		return nil
	}

	platformLabel, err := platform.Format(p, v)
	if err != nil {
		return fmt.Errorf("%w: %w", platform.ErrUnrecognizedPlatform, err)
	}

	log.Infof("installing version: %s (%s)", v, platformLabel)

	if err := install.Executable(ctx, storePath, ex); err != nil {
		return err
	}

	log.Infof("successfully installed version: %s (%s,%s)", ex.Version(), p.OS, p.Arch)

	return nil
}

/* -------------------- Function: resolveVersionFromInput ------------------- */

// Parses command arguments and environment variables and reads pin files to
// determine the correct version of Godot to use.
//
// There are four distinct situations which need to be handled. These are listed
// below along with their desired resolution:
//  1. An explicit version is passed in (e.g. `install <version>`)
//     > The explicitly specified version is returned.
//  2. No version is passed, but the '-g' flag is used (e.g. `install -g`)
//     > The globally pinned version is returned.
//  3. No version is passed, but the '-p' flag is used (e.g. `install --path <path>`)
//     > The pinned version is resolved _at the provided path_.
//  4. No version is passed and no other flag is passed (e.g. `install`)
//     > The pinned version is resolved _in the current working directory_.
//
// For 3. and 4., both of which require pin resolution, the standard resolution
// strategy of checking for a local pin and then a global is used.
func resolveVersionFromInput(c *cli.Context) (version.Version, error) {
	versionArg := c.Args().First()

	v, err := version.Parse(versionArg)
	if err != nil && versionArg != "" {
		return version.Version{}, UsageError{ctx: c, err: err}
	}

	if err == nil {
		return v, nil
	}

	storePath, err := store.Path()
	if err != nil {
		return version.Version{}, err
	}

	// If '-g' is passed then _only_ the globally-pinned version should be
	// returned. Prior validation should have already ensured '-p' was not
	// simultaneously set.
	if c.IsSet("global") && c.Bool("global") {
		return pin.Read(storePath)
	}

	// NOTE: 'filepath.Clean' will replace '' with '.', handling cases 3. and 4.
	// simultaneously.
	path := filepath.Clean(c.String("path"))

	v, err = pin.VersionAt(c.Context, storePath, path) // Update 'v' value.
	if err != nil {
		// Return an error that communicates the root problem and hides any
		// attempted global pin resolution.
		if errors.Is(err, pin.ErrMissingPath) || errors.Is(err, pin.ErrMissingPin) {
			return version.Version{}, fmt.Errorf("%w: %s", pin.ErrMissingPin, path)
		}

		return version.Version{}, err
	}

	return v, nil
}

/* -------------------------- Function: touchStore -------------------------- */

// touchStore determines the store path and ensures it has the expected layout.
func touchStore() (string, error) {
	// Determine the store path.
	storePath, err := store.Path()
	if err != nil {
		return "", err
	}

	// Ensure the store exists.
	if err := store.Touch(storePath); err != nil {
		return "", err
	}

	return storePath, nil
}
