package commands

import (
	"github.com/urfave/cli/v2"
)

// A 'urfave/cli' command to remove a version pin.
func NewUnpin() *cli.Command {
	return &cli.Command{
		Name:  "unpin",
		Usage: "remove a Godot version pin globally or from the specified directory",
		Action: func(c *cli.Context) error {
			return nil
		},
	}
}
