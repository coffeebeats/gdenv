package commands

import (
	"github.com/urfave/cli/v2"
)

// A 'urfave/cli' command to download and cache a specific version of Godot.
func NewInstall() *cli.Command {
	return &cli.Command{
		Name:  "install",
		Usage: "download and cache a specific version of Godot",
		Action: func(c *cli.Context) error {
			return nil
		},
	}
}
