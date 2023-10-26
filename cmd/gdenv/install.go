package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/internal/godot/platform"
	"github.com/coffeebeats/gdenv/internal/godot/version"
	"github.com/coffeebeats/gdenv/pkg/install"
	"github.com/coffeebeats/gdenv/pkg/pin"
	"github.com/coffeebeats/gdenv/pkg/store"
	"github.com/urfave/cli/v2"
)

var ErrInstallUsageGlobalAndPath = errors.New("cannot specify both '-g/--global' and '-p/--path'")

// A 'urfave/cli' command to download and cache a specific version of Godot.
func NewInstall() *cli.Command {
	return &cli.Command{
		Name:     "install",
		Category: "Install",

		Aliases: []string{"i"},

		Usage:     "download and cache a specific version of Godot",
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
				Usage:   "pin the system version",
			},
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Usage:   "determine the version from the pinned `PATH` (cannot be used  with '-g')",
			},
		},

		Action: func(c *cli.Context) error {
			// Validate flag options.
			if c.IsSet("global") && c.IsSet("path") {
				return UsageError{ctx: c, err: ErrPinUsageGlobalAndPath}
			}

			v, err := resolveVersionFromArgOrPath(c)
			if err != nil {
				return err
			}

			if err := installExecutable(c.Context, v, c.Bool("force")); err != nil {
				return err
			}

			if !c.Bool("global") {
				return nil
			}

			// Determine the store path.
			storePath, err := store.Path()
			if err != nil {
				return err
			}

			return writePin(storePath, v)
		},
	}
}

/* ----------------------- Function: installExecutable ---------------------- */

// Installs the specified executable version to the store, but only if needed.
func installExecutable(ctx context.Context, v version.Version, force bool) error {
	log.Infof("installing version: %s", v)

	// Determine the store path.
	storePath, err := store.Path()
	if err != nil {
		return err
	}

	log.Debugf("using store at path: %s", storePath)

	// Ensure the store exists.
	if err := store.Touch(storePath); err != nil {
		return err
	}

	// Define the host 'Platform'.
	p, err := platform.Detect()
	if err != nil {
		return err
	}

	platformLabel, err := platform.Format(p, v)
	if err != nil {
		return fmt.Errorf("%w: %w", platform.ErrUnrecognizedPlatform, err)
	}

	log.Debugf("installing for platform: %s", platformLabel)

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

	if err := install.Executable(ctx, storePath, ex); err != nil {
		return err
	}

	log.Infof("successfully installed version: %s (%s,%s)", ex.Version(), p.OS, p.Arch)

	return nil
}

/* ------------------ Function: resolveVersionFromArgOrPath ----------------- */

func resolveVersionFromArgOrPath(c *cli.Context) (version.Version, error) {
	versionArg := c.Args().First()

	v, err := version.Parse(versionArg)
	if err != nil && versionArg != "" {
		return version.Version{}, UsageError{ctx: c, err: err}
	}

	if err == nil {
		return v, nil
	}

	path := c.String("path")
	if path == "" {
		path, err = os.Getwd() // Update 'path' value.
		if err != nil {
			return version.Version{}, err
		}
	}

	// NOTE: Omit store path to avoid resolving the global pin version.
	v, err = pin.VersionAt(c.Context, "", path) // Update 'v' value.
	if err != nil {
		// Return an error that communicates the root problem and hides the
		// storePath="" hack from above.
		if errors.Is(err, pin.ErrMissingPath) {
			return version.Version{}, fmt.Errorf("%w: %s", pin.ErrMissingPin, path)
		}

		return version.Version{}, err
	}

	return v, nil
}
