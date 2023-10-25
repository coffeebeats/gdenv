package main

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/internal/godot/platform"
	"github.com/coffeebeats/gdenv/internal/godot/version"
	"github.com/coffeebeats/gdenv/pkg/install"
	"github.com/coffeebeats/gdenv/pkg/pin"
	"github.com/coffeebeats/gdenv/pkg/store"
	"github.com/urfave/cli/v2"
)

// A 'urfave/cli' command to download and cache a specific version of the Godot
// source code.
func NewVendor() *cli.Command {
	return &cli.Command{
		Name:     "vendor",
		Category: "Vendor",

		Usage:     "download and cache a specific version of Godot source code",
		UsageText: "gdenv install [OPTIONS] <VERSION>",

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
		},

		Action: func(c *cli.Context) error {
			// Validate arguments
			v, err := version.Parse(c.Args().First())
			if err != nil {
				return UsageError{ctx: c, err: err}
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

			if err := pin.Write(v, storePath); err != nil {
				return err
			}

			log.Infof("set system default version: %s", v)

			return nil
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
