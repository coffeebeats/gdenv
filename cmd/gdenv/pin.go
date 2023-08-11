package main

import (
	"github.com/urfave/cli/v2"
)

// A 'urfave/cli' command to pin a Godot version globally or for a directory.
func NewPin() *cli.Command {
	return &cli.Command{
		Name:     "pin",
		Category: "Pin",

		Usage:     "set the Godot version globally or for a specific directory",
		UsageText: "gdenv pin [OPTIONS] <VERSION>",

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "global",
				Aliases: []string{"g"},
				Usage:   "pin the system version (cannot be used with '-p')",
			},
			&cli.BoolFlag{
				Name:    "install",
				Aliases: []string{"i"},
				Usage:   "installs the specified version of Godot if missing",
			},
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Usage:   "pin the specified `PATH` (cannot be used with '-g')",
			},
		},

		Action: pin,
	}
}

func pin(_ *cli.Context) error {
	return nil
}
