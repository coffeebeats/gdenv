package main

import (
	"log"

	"github.com/coffeebeats/gdenv/internal/godot/platform"
	"github.com/coffeebeats/gdenv/pkg/install"
	"github.com/coffeebeats/gdenv/pkg/store"
	"github.com/urfave/cli/v2"
)

// A 'urfave/cli' command to print the path to the effective Godot binary.
func NewWhich() *cli.Command {
	return &cli.Command{
		Name:     "which",
		Category: "Utilities",

		Usage:     "print the path to the Godot executable which would be used in the specified directory",
		UsageText: "gdenv which [OPTIONS]",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Usage:   "Check the specified `PATH`",
			},
		},

		Action: func(c *cli.Context) error {
			// Determine 'path' option
			path, err := resolvePath(c)
			if err != nil {
				return err
			}

			// Ensure 'Store' layout
			storePath, err := store.InitAtPath()
			if err != nil {
				return err
			}

			// Define the host 'Platform'.
			platform, err := platform.Detect()
			if err != nil {
				return err
			}

			toolPath, err := install.Which(c.Context, storePath, platform, path)
			if err != nil {
				return err
			}

			log.Println(toolPath)

			return nil
		},
	}
}
