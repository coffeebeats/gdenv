package main

import (
	"fmt"
	"log"
	"os"

	"github.com/coffeebeats/gdenv/internal/commands"
	"github.com/urfave/cli/v2"
)

func main() {
	cli.VersionPrinter = versionPrinter

	app := &cli.App{
		Name:    "gdenv",
		Version: "v0.1.1", // x-release-please-version

		Commands: []*cli.Command{

			/* -------------------------------- Pin/Unpin ------------------------------- */

			commands.NewPin(),
			commands.NewUnpin(),

			/* ---------------------------- Install/Uninstall --------------------------- */

			commands.NewInstall(),
			commands.NewUninstall(),

			/* --------------------------------- Utility -------------------------------- */

			commands.NewCompletions(),
			commands.NewLs(),
			commands.NewWhich(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func versionPrinter(cCtx *cli.Context) {
	fmt.Printf("gdenv %s\n", cCtx.App.Version)
}
