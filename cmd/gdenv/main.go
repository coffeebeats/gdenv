package main

import (
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

func fail(err error) cli.ExitCoder {
	return cli.Exit(err, 1)
}

func versionPrinter(cCtx *cli.Context) {
	log.Printf("gdenv %s\n", cCtx.App.Version)
}
