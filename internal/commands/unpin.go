package commands

import (
	"os"

	"github.com/coffeebeats/gdenv/pkg/pin"
	"github.com/coffeebeats/gdenv/pkg/store"
	"github.com/urfave/cli/v2"
)

// A 'urfave/cli' command to remove a version pin.
func NewUnpin() *cli.Command {
	return &cli.Command{
		Name:     "unpin",
		Category: "Pin",

		Usage:     "remove a Godot version pin globally or from the specified directory",
		UsageText: "gdenv unpin [OPTIONS]",

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "global",
				Aliases: []string{"g"},
				Usage:   "unpin the system version (cannot be used with '-p')",
			},
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Usage:   "unpin the specified `PATH` (cannot be used with '-g')",
			},
		},

		Action: func(c *cli.Context) error {
			// Validate flag options.
			if c.IsSet("global") && c.IsSet("path") {
				return cli.Exit("cannot specify both '--global' and '--path'", 1)
			}

			// Ensure 'Store' layout
			s, err := store.Path()
			if err != nil {
				return cli.Exit(err, 1)
			}

			if err := store.Init(s); err != nil {
				return cli.Exit(err, 1)
			}

			// Determine 'path' option
			var path string
			switch {
			case c.Bool("path"):
				path = c.String("path")
			case c.Bool("global"):
				p, err := store.Path()
				if err != nil {
					return cli.Exit(err, 1)
				}

				path = p
			default:
				p, err := os.Getwd()
				if err != nil {
					return cli.Exit(err, 1)
				}

				path = p
			}

			return pin.Remove(path)
		},
	}
}
