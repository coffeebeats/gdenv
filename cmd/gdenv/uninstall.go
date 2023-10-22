package main

import (
	"github.com/charmbracelet/log"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/internal/godot/platform"
	"github.com/coffeebeats/gdenv/internal/godot/version"
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
			// Determine the store path.
			storePath, err := store.Path()
			if err != nil {
				return err
			}

			// Ensure the store's layout is correct.
			if err := store.Touch(storePath); err != nil {
				return err
			}

			// Define the host 'Platform'.
			p, err := platform.Detect()
			if err != nil {
				return err
			}

			// Uninstall all versions.
			if c.Bool("all") {
				ee, err := store.Executables(c.Context, storePath)
				if err != nil {
					return err
				}

				if len(ee) == 0 {
					return nil
				}

				log.Info("removing all installed versions")

				return store.Clear(storePath)
			}

			// Uninstall a specific version.

			// Validate arguments
			v, err := version.Parse(c.Args().First())
			if err != nil && !c.Bool("all") {
				return UsageError{ctx: c, err: err}
			}

			// Define the target 'Executable'.
			ex := executable.New(v, p)

			ok, err := store.Has(storePath, ex)
			if err != nil {
				return err
			}

			if !ok {
				return nil
			}

			log.Infof("uninstalling version: %s", v)

			return store.Remove(storePath, ex)
		},
	}
}
