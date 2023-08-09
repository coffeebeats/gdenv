package commands

import (
	"github.com/urfave/cli/v2"
)

// A 'urfave/cli' command to print the path to the effective Godot binary.
func NewWhich() *cli.Command {
	return &cli.Command{
		Name:  "which",
		Usage: "print the path to the Godot executable which would be used in the specified directory",

		Action: which,
	}
}

func which(c *cli.Context) error {
	return nil
}
