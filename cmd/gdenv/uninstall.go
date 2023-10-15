package main

import (
	"github.com/coffeebeats/gdenv/internal/godot/artifact/executable"
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
			// Ensure 'Store' layout
			storePath, err := store.InitAtPath()
			if err != nil {
				return err
			}

			// Define the host 'Platform'.
			p, err := detectPlatform()
			if err != nil {
				return err
			}

			// Uninstall a specific version.
			if !c.Bool("all") {
				// Validate arguments
				v, err := version.Parse(c.Args().First())
				if err != nil && !c.Bool("all") {
					return UsageError{ctx: c, err: err}
				}

				// Define the target 'Executable'.
				ex := executable.New(v, p)

				return uninstall(storePath, ex)
			}

			// Uninstall all versions.
			return uninstallAll(storePath)
		},
	}
}

/* --------------------------- Function: uninstall -------------------------- */

// Deletes a platform-specific version of Godot from the store.
func uninstall(storePath string, ex executable.Executable) error {
	return store.Remove(storePath, ex)
}

/* ------------------------- Function: uninstallAll ------------------------- */

// Uninstalls all cached Godot executables, regardless of platform.
func uninstallAll(storePath string) error {
	executables, err := store.Executables(storePath)
	if err != nil {
		return err
	}

	for _, ex := range executables {
		if err := uninstall(storePath, ex); err != nil {
			return err
		}
	}

	return nil
}
