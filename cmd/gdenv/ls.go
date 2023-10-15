package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/coffeebeats/gdenv/pkg/pin"
	"github.com/coffeebeats/gdenv/pkg/store"
	"github.com/urfave/cli/v2"
)

/* ----------------------------- Function: NewLs ---------------------------- */

// A 'urfave/cli' command to print installed versions of Godot.
func NewLs() *cli.Command {
	return &cli.Command{
		Name:     "ls",
		Category: "Utilities",

		Usage:     "print the path and version of all of the installed versions of Godot",
		UsageText: "gdenv ls",

		Action: func(c *cli.Context) error {
			// Ensure 'Store' layout
			storePath, err := store.InitAtPath()
			if err != nil {
				return err
			}

			results, err := ls(c.Context, storePath)
			if err != nil {
				return err
			}

			if wd, err := os.Getwd(); err == nil {
				if version, err := pin.Read(wd); err == nil {
					log.Printf("🤖 Currently active version: %s\n\n", version)
				}
			}

			log.Printf("Installed versions:\n\n%s\n", strings.Join(results, "\n"))

			return nil
		},
	}
}

/* ------------------------------ Function: ls ------------------------------ */

func ls(ctx context.Context, storePath string) ([]string, error) {
	executables, err := store.Executables(ctx, storePath)
	if err != nil {
		return nil, err
	}

	out := make([]string, len(executables))

	for i, ex := range executables {
		select {
		case <-ctx.Done():
			return out, ctx.Err()
		default:
		}

		out[i] = ex.String()
	}

	return out, nil
}
