package commands

import (
	"github.com/urfave/cli/v2"
)

// A 'urfave/cli' command to remove a version pin.
func NewUnpin() *cli.Command {
	return &cli.Command{
		Name:     "unpin",
		Category: "Pin",

		Usage:     "remove a Godot version pin globally or from the specified directory",
		UsageText: "gdenv unpin [OPTIONS]",

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "global",
				Aliases: []string{"g"},
				Usage:   "unpin the system version (cannot be used with '-p')",
			},
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Usage:   "unpin the specified `PATH` (cannot be used with '-g')",
			},
		},

		Action: unpin,
	}
}

func unpin(_ *cli.Context) error {
	return nil
}
