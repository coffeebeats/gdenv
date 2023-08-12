package main

import (
	"github.com/coffeebeats/gdenv/internal/godot"
	"github.com/coffeebeats/gdenv/pkg/store"
	"github.com/urfave/cli/v2"
)

// A 'urfave/cli' command to delete a cached version of Godot.
func NewUninstall() *cli.Command {
	return &cli.Command{
		Name:     "uninstall",
		Category: "Install",

		Usage:     "remove the specified version of Godot from the gdenv download cache",
		UsageText: "gdenv uninstall [OPTIONS] [VERSION]",

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "uninstall all versions of Godot in the cache",
			},
		},

		Action: func(c *cli.Context) error {
			// Ensure 'Store' layout
			storePath, err := store.InitAtPath()
			if err != nil {
				return err
			}

			// Uninstall a specific version.
			if !c.Bool("all") {
				// Validate arguments
				version, err := godot.ParseVersion(c.Args().First())
				if err != nil && !c.Bool("all") {
					return fail(err)
				}

				if err := store.Remove(storePath, version); err != nil {
					return fail(err)
				}

				return nil
			}

			// Uninstall all versions.
			versions, err := store.Versions(storePath)
			if err != nil {
				return fail(err)
			}

			for _, v := range versions {
				if err := store.Remove(storePath, v); err != nil {
					return fail(err)
				}
			}

			return nil
		},
	}
}
