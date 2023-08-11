package main

import (
	"github.com/urfave/cli/v2"
)

// A 'urfave/cli' command to print installed versions of Godot.
func NewLs() *cli.Command {
	return &cli.Command{
		Name:     "ls",
		Category: "Utilities",

		Usage:     "print the path and version of all of the installed versions of Godot",
		UsageText: "gdenv ls",

		Action: ls,
	}
}

func ls(_ *cli.Context) error {
	return nil
}
