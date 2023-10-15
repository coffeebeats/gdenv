package main

import (
	"context"

	"github.com/coffeebeats/gdenv/internal/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/internal/godot/platform"
	"github.com/coffeebeats/gdenv/internal/godot/version"
	"github.com/coffeebeats/gdenv/pkg/store"
	"github.com/urfave/cli/v2"
)

const (
	EnvGDEnvArch     = "GDENV_ARCH"
	EnvGDEnvOS       = "GDENV_OS"
	EnvGDEnvPlatform = "GDENV_PLATFORM"
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

			// Ensure 'Store' layout
			storePath, err := store.InitAtPath()
			if err != nil {
				return err
			}

			// Define the host 'Platform'.
			p, err := platform.Detect()
			if err != nil {
				return err
			}

			// Define the target 'Executable'.
			ex := executable.New(v, p)

			if store.Has(storePath, ex) && !c.Bool("force") {
				return nil
			}

			return install(c.Context, storePath, ex)
		},
	}
}

/* ---------------------------- Function: install --------------------------- */

// Downloads and caches a platform-specific version of Godot.
func install(_ context.Context, _ string, _ executable.Executable) error {
	return nil
}
