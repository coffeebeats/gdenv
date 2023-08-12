package main

import (
	"github.com/coffeebeats/gdenv/internal/godot"
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
				return fail(err)
			}

			// Ensure 'Store' layout
			storePath, err := store.InitAtPath()
			if err != nil {
				return fail(err)
			}

			if store.Has(storePath, version) && !c.Bool("force") {
				return nil
			}

			if err := install(storePath, version); err != nil {
				return fail(err)
			}

			return nil
		},
	}
}

/* ---------------------------- Function: install --------------------------- */

func install(_ string, _ godot.Version) error {
	return nil
}
