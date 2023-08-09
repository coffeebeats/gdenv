package commands

import (
	"github.com/urfave/cli/v2"
)

var Ls = cli.Command{
	Name:  "ls",
	Usage: "print the path and version of all of the installed versions of Godot",
	Action: func(cCtx *cli.Context) error {
		return nil
	},
}
