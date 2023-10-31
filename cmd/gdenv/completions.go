package main

import (
	"github.com/urfave/cli/v2"
)

// A 'urfave/cli' command to print shell completions.
func NewCompletions() *cli.Command {
	return &cli.Command{
		Name:     "completions",
		Category: "Utilities",

		Usage:     "print shell completions for the 'gdenv' CLI application",
		UsageText: "gdenv completions [OPTIONS] <SHELL>",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"p"},
				Usage:   "file to write the completions to",
			},
		},

		Action: func(c *cli.Context) error {
			return nil
		},
	}
}
