package main

import (
	"github.com/coffeebeats/gdenv/pkg/godot"
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
				return fail(err)
			}

			// Define the host 'Platform'.
			platform, err := godot.HostPlatform()
			if err != nil {
				return fail(err)
			}

			// Uninstall a specific version.
			if !c.Bool("all") {
				// Validate arguments
				version, err := godot.ParseVersion(c.Args().First())
				if err != nil && !c.Bool("all") {
					return failWithUsage(c, err)
				}

				// Define the target 'Executable'.
				ex := godot.Executable{Platform: platform, Version: version}

				if err := uninstall(storePath, ex); err != nil {
					return fail(err)
				}

				return nil
			}

			// Uninstall all versions.
			if err := uninstallAll(storePath); err != nil {
				return fail(err)
			}

			return nil
		},
	}
}

/* --------------------------- Function: uninstall -------------------------- */

// Deletes a platform-specific version of Godot from the store.
func uninstall(storePath string, ex godot.Executable) error {
	if err := store.Remove(storePath, ex); err != nil {
		return fail(err)
	}

	return nil
}

/* ------------------------- Function: uninstallAll ------------------------- */

// Uninstalls all cached Godot executables, regardless of platform.
func uninstallAll(storePath string) error {
	executables, err := store.Executables(storePath)
	if err != nil {
		return fail(err)
	}

	for _, ex := range executables {
		if err := uninstall(storePath, ex); err != nil {
			return fail(err)
		}
	}

	return nil
}
