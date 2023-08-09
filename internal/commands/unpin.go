package commands

import (
	"github.com/urfave/cli/v2"
)

// A 'urfave/cli' command to remove a version pin.
func NewUnpin() *cli.Command {
	return &cli.Command{
		Name:      "unpin",
		Usage:     "remove a Godot version pin globally or from the specified directory",
		UsageText: "gdenv pin [OPTIONS] <VERSION>",

		Action: unpin,
	}
}

func unpin(c *cli.Context) error {
	return nil
}
