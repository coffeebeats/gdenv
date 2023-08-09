package commands

import (
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

		Action: uninstall,
	}
}

func uninstall(_ *cli.Context) error {
	return nil
}
