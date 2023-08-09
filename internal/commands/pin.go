package commands

import (
	"github.com/urfave/cli/v2"
)

// A 'urfave/cli' command to pin a Godot version globally or for a directory.
func NewPin() *cli.Command {
	return &cli.Command{
		Name:  "pin",
		Usage: "set the Godot version globally or for a specific directory",
		Action: func(c *cli.Context) error {
			return nil
		},
	}
}
