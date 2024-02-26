package main

import (
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v2"

	"github.com/coffeebeats/gdenv/pkg/install"
)

// A 'urfave/cli' command to download and cache a specific version of the Godot
// source code and extract it into the specified directory.
func NewVendor() *cli.Command {
	return &cli.Command{
		Name:     "vendor",
		Category: "Install",

		Usage:     "download the Godot source code to the specified directory",
		UsageText: "gdenv install [OPTIONS] [VERSION]",

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "forcibly overwrite an existing cache entry",
			},
			&cli.StringFlag{
				Name:    "out",
				Aliases: []string{"o"},
				Value:   "./godot",
				Usage:   "extract the source code into 'OUT' (overwrites conflicting files)",
			},
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Usage:   "resolve the pinned 'VERSION' at 'PATH'",
			},
		},

		Action: func(c *cli.Context) error {
			v, err := resolveVersionFromInput(c)
			if err != nil {
				return err
			}

			storePath, err := touchStore()
			if err != nil {
				return err
			}

			log.Debugf("using store at path: %s", storePath)

			if err := installSource(c.Context, storePath, v, c.Bool("force")); err != nil {
				return err
			}

			return install.Vendor(c.Context, v, storePath, c.String("out"))
		},
	}
}
