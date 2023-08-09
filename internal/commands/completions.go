package commands

import (
	"github.com/urfave/cli/v2"
)

// A 'urfave/cli' command to print shell completions.
func NewCompletions() *cli.Command {
	return &cli.Command{
		Name:  "completions",
		Usage: "print shell completions for the gdenv CLI application",
		Action: func(c *cli.Context) error {
			return nil
		},
	}
}
