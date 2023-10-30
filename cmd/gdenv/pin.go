package main

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/coffeebeats/gdenv/pkg/godot/version"
	"github.com/coffeebeats/gdenv/pkg/pin"
	"github.com/urfave/cli/v2"
)

var (
	ErrPinUsageForceAndInstall = errors.New("cannot specify '-f/--force' without '-i/--install'")
	ErrPinUsageGlobalAndPath   = errors.New("cannot specify both '-g/--global' and '-p/--path'")
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
				return UsageError{ctx: c, err: ErrPinUsageGlobalAndPath}
			}

			if c.IsSet("force") && !c.IsSet("install") {
				return UsageError{ctx: c, err: ErrPinUsageForceAndInstall}
			}

			// Validate arguments
			v, err := version.Parse(c.Args().First())
			if err != nil {
				return UsageError{ctx: c, err: err}
			}

			storePath, err := touchStore()
			if err != nil {
				return err
			}

			// Determine 'path' option
			pinPath, err := resolvePath(c)
			if err != nil {
				return err
			}

			if err := writePin(storePath, pinPath, v); err != nil {
				return err
			}

			if !c.Bool("install") {
				return nil
			}

			return installExecutable(c.Context, storePath, v, c.Bool("force"))
		},
	}
}

/* -------------------------- Function: resolvePath ------------------------- */

// Determines the path to pin based on the provided options.
func resolvePath(c *cli.Context) (string, error) {
	switch {
	case c.IsSet("path"):
		return filepath.Clean(c.String("path")), nil
	case c.Bool("global"):
		return touchStore()
	default:
		p, err := os.Getwd()
		if err != nil {
			return "", err
		}

		return p, nil
	}
}

/* --------------------------- Function: writePin --------------------------- */

// Writes the specified version to a pin file.
func writePin(storePath, pinPath string, v version.Version) error {
	if err := pin.Write(v, pinPath); err != nil {
		return err
	}

	if pinPath == storePath {
		log.Infof("set system default version: %s", v)
	} else {
		log.Infof("pinned '%s' to version: %s", pinPath, v)
	}

	return nil
}
