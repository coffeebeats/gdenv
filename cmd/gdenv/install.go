package main

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/internal/godot/platform"
	"github.com/coffeebeats/gdenv/internal/godot/version"
	"github.com/coffeebeats/gdenv/pkg/install"
	"github.com/coffeebeats/gdenv/pkg/store"
	"github.com/urfave/cli/v2"
)

// A 'urfave/cli' command to download and cache a specific version of Godot.
func NewInstall() *cli.Command {
	return &cli.Command{
		Name:     "install",
		Category: "Install",

		Usage:     "download and cache a specific version of Godot",
		UsageText: "gdenv install [OPTIONS] <VERSION>",

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "forcibly overwrite an existing cache entry",
			},
		},

		Action: func(c *cli.Context) error {
			// Validate arguments
			v, err := version.Parse(c.Args().First())
			if err != nil {
				return UsageError{ctx: c, err: err}
			}

			log.Infof("installing version: %s", v)

			return installExecutable(c.Context, v, c.Bool("force"))
		},
	}
}

/* ----------------------- Function: installExecutable ---------------------- */

// Installs the specified executable version to the store, but only if needed.
func installExecutable(ctx context.Context, v version.Version, force bool) error {
	// Determine the store path.
	storePath, err := store.Path()
	if err != nil {
		return err
	}

	log.Debugf("using store at path: %s", storePath)

	// Ensure the store's layout is correct.
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

	path, err := store.Executable(storePath, ex)
	if err != nil {
		return err
	}

	log.Infof("successfully installed version: %s", path)

	return nil
}
