package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	log.SetFlags(0)

	cli.VersionPrinter = versionPrinter

	app := &cli.App{
		Name:    "gdenv",
		Version: "v0.1.1", // x-release-please-version

		Suggest:                true,
		UseShortOptionHandling: true,

		Commands: []*cli.Command{

			/* -------------------------------- Pin/Unpin ------------------------------- */

			NewPin(),
			NewUnpin(),

			/* ---------------------------- Install/Uninstall --------------------------- */

			NewInstall(),
			NewUninstall(),

			/* --------------------------------- Utility -------------------------------- */

			NewCompletions(),
			NewLs(),
			NewWhich(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// Convenience function which return an error that invokes 'os.Exit(1)'.
func fail(err error) error {
	return cli.Exit(err, 1)
}

func failWithUsage(c *cli.Context, err error) error {
	cli.ShowAppHelp(c)
	log.Println()

	return cli.Exit(fmt.Errorf("command failed: %w", err), 1)
}

func versionPrinter(cCtx *cli.Context) {
	log.Printf("gdenv %s\n", cCtx.App.Version)
}
