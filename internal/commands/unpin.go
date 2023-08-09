package commands

import (
	"github.com/urfave/cli/v2"
)

var Unpin = cli.Command{
	Name:  "unpin",
	Usage: "remove a Godot version pin globally or from the specified directory",
	Action: func(cCtx *cli.Context) error {
		return nil
	},
}
