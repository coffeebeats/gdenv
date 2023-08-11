package commands

import (
	"os"

	"github.com/coffeebeats/gdenv/internal/godot"
	"github.com/coffeebeats/gdenv/pkg/pin"
	"github.com/coffeebeats/gdenv/pkg/store"
	"github.com/urfave/cli/v2"
)

// A 'urfave/cli' command to pin a Godot version globally or for a directory.
func NewPin() *cli.Command {
	return &cli.Command{
		Name:     "pin",
		Category: "Pin",

		Usage:     "set the Godot version globally or for a specific directory",
		UsageText: "gdenv pin [OPTIONS] <VERSION>",

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "global",
				Aliases: []string{"g"},
				Usage:   "pin the system version (cannot be used with '-p')",
			},
			&cli.BoolFlag{
				Name:    "install",
				Aliases: []string{"i"},
				Usage:   "installs the specified version of Godot if missing",
			},
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Usage:   "pin the specified `PATH` (cannot be used with '-g')",
			},
		},

		Action: func(c *cli.Context) error {
			// Validate flag options.
			if c.IsSet("global") && c.IsSet("path") {
				return cli.Exit("cannot specify both '--global' and '--path'", 1)
			}

			// Validate arguments
			version, err := godot.ParseVersion(c.Args().First())
			if err != nil {
				return cli.Exit(err, 1)
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

			if c.Bool("install") {
				install(version)
			}

			return pin.Write(version, path)
		},
	}
}
