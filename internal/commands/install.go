package commands

import (
	"github.com/urfave/cli/v2"
)

var Install = cli.Command{
	Name:  "install",
	Usage: "download and cache a specific version of Godot",
	Action: func(cCtx *cli.Context) error {
		return nil
	},
}
