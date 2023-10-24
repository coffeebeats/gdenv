package main

import (
	"errors"
	"os"

	"github.com/charmbracelet/log"
	"github.com/coffeebeats/gdenv/internal/godot/version"
	"github.com/coffeebeats/gdenv/pkg/pin"
	"github.com/coffeebeats/gdenv/pkg/store"
	"github.com/urfave/cli/v2"
)

var (
	ErrMissingPin           = errors.New("missing version pin")
	ErrUsageForceAndInstall = errors.New("cannot specify '-f/--force' without '-i/--install'")
	ErrUsageGlobalAndPath   = errors.New("cannot specify both '-g/--global' and '-p/--path'")
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
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "forcibly overwrite an existing cache entry (only used with '-i')",
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
				return UsageError{ctx: c, err: ErrUsageGlobalAndPath}
			}

			if c.IsSet("force") && !c.IsSet("install") {
				return UsageError{ctx: c, err: ErrUsageForceAndInstall}
			}

			// Validate arguments
			v, err := version.Parse(c.Args().First())
			if err != nil {
				return UsageError{ctx: c, err: err}
			}

			// Determine 'path' option
			pinPath, err := resolvePath(c)
			if err != nil {
				return err
			}

			if err := pin.Write(v, pinPath); err != nil {
				return err
			}

			// Determine the store path.
			storePath, err := store.Path()
			if err != nil {
				return err
			}

			if pinPath == storePath {
				log.Infof("set system default version: %s", v)
			} else {
				log.Infof("pinned '%s' to version: %s", pinPath, v)
			}

			if !c.Bool("install") {
				return nil
			}

			return installExecutable(c.Context, v, c.Bool("force"))
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
