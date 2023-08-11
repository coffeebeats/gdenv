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
				return cli.Exit(err, 1)
			}

			// Ensure 'Store' layout
			s, err := store.Path()
			if err != nil {
				return cli.Exit(err, 1)
			}

			if err := store.Init(s); err != nil {
				return cli.Exit(err, 1)
			}

			if store.Has(s, version) && !c.Bool("force") {
				return nil
			}

			if err := install(version); err != nil {
				return cli.Exit(err, 1)
			}

			return nil
		},
	}
}

func install(version godot.Version) error {
	return nil
}
