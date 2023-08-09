package commands

import (
	"github.com/urfave/cli/v2"
)

var Which = cli.Command{
	Name:  "which",
	Usage: "print the path to the Godot executable which would be used in the specified directory",
	Action: func(cCtx *cli.Context) error {
		return nil
	},
}
