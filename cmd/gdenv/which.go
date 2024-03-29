package main

import (
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v2"

	"github.com/coffeebeats/gdenv/pkg/godot/platform"
	"github.com/coffeebeats/gdenv/pkg/install"
)

// A 'urfave/cli' command to print the path to the effective Godot binary.
func NewWhich() *cli.Command {
	return &cli.Command{
		Name:     "which",
		Category: "Utilities",

		Usage:     "print the path to the Godot executable which would be used in the specified directory",
		UsageText: "gdenv which [OPTIONS]",

		Flags: []cli.Flag{
			newVerboseFlag(),

			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Usage:   "check at the specified `PATH`",
			},
		},

		Action: func(c *cli.Context) error {
			// Determine 'path' option
			pinPath, err := resolvePath(c)
			if err != nil {
				return err
			}

			// Determine the store path.
			storePath, err := touchStore()
			if err != nil {
				return err
			}

			log.Debugf("using store at path: %s", storePath)

			// Define the host 'Platform'.
			p, err := platform.Detect()
			if err != nil {
				return err
			}

			path, err := install.Which(c.Context, storePath, p, pinPath)
			if err != nil {
				return err
			}

			if path != "" {
				log.Print(path)
			}

			return nil
		},
	}
}
