package commands

import (
	"github.com/urfave/cli/v2"
)

var Uninstall = cli.Command{
	Name:  "uninstall",
	Usage: "remove the specified version of Godot from the gdenv download cache",
	Action: func(cCtx *cli.Context) error {
		return nil
	},
}
