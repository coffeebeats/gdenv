package commands

import (
	"github.com/urfave/cli/v2"
)

var Pin = cli.Command{
	Name:  "pin",
	Usage: "set the Godot version globally or for a specific directory",
	Action: func(cCtx *cli.Context) error {
		return nil
	},
}
