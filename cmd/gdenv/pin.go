package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/coffeebeats/gdenv/pkg/godot"
	"github.com/coffeebeats/gdenv/pkg/pin"
	"github.com/coffeebeats/gdenv/pkg/store"
	"github.com/urfave/cli/v2"
)

var (
	ErrOptionUsage = errors.New("gdenv: invalid option usage")
)

/* ---------------------------- Function: NewPin ---------------------------- */

// A 'urfave/cli' command to pin a Godot version globally or for a directory.
func NewPin() *cli.Command { //nolint:funlen
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
				err := fmt.Errorf("%w: cannot specify both '--global' and '--path'", ErrOptionUsage)
				return failWithUsage(c, err)
			}

			// Determine 'path' option
			path, err := resolvePath(c)
			if err != nil {
				return fail(err)
			}

			// Validate arguments
			version, err := godot.ParseVersion(c.Args().First())
			if err != nil {
				return failWithUsage(c, err)
			}

			if c.Bool("install") {
				// Ensure 'Store' layout
				storePath, err := store.InitAtPath()
				if err != nil {
					return fail(err)
				}

				// Define the host 'Platform'.
				platform, err := detectPlatform()
				if err != nil {
					return fail(err)
				}

				// Define the target 'Executable'.
				ex := godot.Executable{Platform: platform, Version: version}

				if err := install(storePath, ex); err != nil {
					return fail(err)
				}
			}

			if err := pin.Write(version, path); err != nil {
				return fail(err)
			}

			return nil
		},
	}
}

/* -------------------------- Function: resolvePath ------------------------- */

// Determines the path to pin based on the provided options.
func resolvePath(c *cli.Context) (string, error) {
	switch {
	case c.IsSet("path"):
		return c.String("path"), nil
	case c.Bool("global"):
		p, err := store.Path()
		if err != nil {
			return "", err
		}

		return p, nil
	default:
		p, err := os.Getwd()
		if err != nil {
			return "", err
		}

		return p, nil
	}
}
