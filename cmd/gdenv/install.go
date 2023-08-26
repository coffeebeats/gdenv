package main

import (
	"github.com/coffeebeats/gdenv/pkg/godot"
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
			version, err := godot.ParseVersion(c.Args().First())
			if err != nil {
				return failWithUsage(c, err)
			}

			// Ensure 'Store' layout
			storePath, err := store.InitAtPath()
			if err != nil {
				return fail(err)
			}

			// Define the host 'Platform'.
			platform, err := godot.HostPlatform()
			if err != nil {
				return fail(err)
			}

			// Define the target 'Executable'.
			ex := godot.Executable{Platform: platform, Version: version}

			if store.Has(storePath, ex) && !c.Bool("force") {
				return nil
			}

			if err := install(storePath, ex); err != nil {
				return fail(err)
			}

			return nil
		},
	}
}

/* ---------------------------- Function: install --------------------------- */

// Downloads and caches a platform-specific version of Godot.
func install(_ string, _ godot.Executable) error {
	return nil
}
