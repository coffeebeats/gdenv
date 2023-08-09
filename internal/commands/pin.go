package commands

import (
	"strings"

	"github.com/urfave/cli/v2"
)

const globalFlag = "global"
const installFlag = "install"
const pathFlag = "path"

// A 'urfave/cli' command to pin a Godot version globally or for a directory.
func NewPin() *cli.Command {
	return &cli.Command{
		Name:      "pin",
		Usage:     "set the Godot version globally or for a specific directory",
		UsageText: "gdenv pin [OPTIONS] <VERSION>",

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    globalFlag,
				Aliases: strings.Split(globalFlag, "")[0:1],
				Usage:   "pin the system version (cannot be used with '-p')",
			},
			&cli.BoolFlag{
				Name:    installFlag,
				Aliases: strings.Split(installFlag, "")[0:1],
				Usage:   "installs the specified version of Godot if missing",
			},
			&cli.StringFlag{
				Name:    pathFlag,
				Aliases: strings.Split(pathFlag, "")[0:1],
				Usage:   "pin the specified `PATH` (cannot be used with '-g')",
			},
		},

		Action: pin,
	}
}

func pin(c *cli.Context) error {
	return nil
}
