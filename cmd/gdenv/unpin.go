package main

import (
	"errors"

	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v2"

	"github.com/coffeebeats/gdenv/pkg/pin"
	"github.com/coffeebeats/gdenv/pkg/store"
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
				Usage:   "unpin the specified 'PATH' (cannot be used with '-g')",
			},
		},

		Action: func(c *cli.Context) error {
			// Validate flag options.
			if c.IsSet("global") && c.IsSet("path") {
				return UsageError{ctx: c, err: ErrPinUsageGlobalAndPath}
			}

			// Determine 'path' option
			pinPath, err := resolvePath(c)
			if err != nil {
				return err
			}

			// Exit early if the pin doesn't exist.
			if _, err := pin.Read(pinPath); err != nil {
				if !errors.Is(err, pin.ErrMissingPin) {
					return err
				}

				return nil
			}

			if err := pin.Remove(pinPath); err != nil {
				return err
			}

			// Determine the store path.
			storePath, err := store.Path()
			if err != nil {
				return err
			}

			if pinPath == storePath {
				log.Info("unset system default version")
			} else {
				log.Infof("removed version pin from path: %s", pinPath)
			}

			return nil
		},
	}
}
