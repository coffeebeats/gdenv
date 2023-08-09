package commands

import (
	"github.com/urfave/cli/v2"
)

// A 'urfave/cli' command to delete a cached version of Godot.
func NewUninstall() *cli.Command {
	return &cli.Command{
		Name:  "uninstall",
		Usage: "remove the specified version of Godot from the gdenv download cache",
		Action: func(c *cli.Context) error {
			return nil
		},
	}
}
