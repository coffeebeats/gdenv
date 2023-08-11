package main

import (
	"github.com/urfave/cli/v2"
)

// A 'urfave/cli' command to print shell completions.
func NewCompletions() *cli.Command {
	return &cli.Command{
		Name:     "completions",
		Category: "Utilities",

		Usage:     "print shell completions for the gdenv CLI application",
		UsageText: "gdenv completions [OPTIONS] <SHELL>",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"p"},
				Usage:   "Write the completions to `OUT_FILE`",
			},
		},

		Action: completions,
	}
}

func completions(_ *cli.Context) error {
	return nil
}
