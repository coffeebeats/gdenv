package commands

import (
	"github.com/urfave/cli/v2"
)

var Completions = cli.Command{
	Name:  "completions",
	Usage: "print shell completions for the gdenv CLI application",
	Action: func(cCtx *cli.Context) error {
		return nil
	},
}
