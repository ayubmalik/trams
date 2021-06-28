package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var version = "dev"

func main() {
	app := &cli.App{
		Usage:     "display Metrolink tram information using data from TfGM API",
		UsageText: "trams COMMAND [args]",
		Commands: []*cli.Command{
			{
				Name:      "display",
				Usage:     "display tram information for specified Metrolink stations. If no stations are specified displays all stations. Run `trams help display` for more details.",
				UsageText: "display [STATION...] - If no STATION arguments are specified, displays all stations. Multiple STATION arguments can be specified as short name or long name, e.g. `display BCH MAN` or `display Benchill MAN`",
				Action: func(c *cli.Context) error {
					fmt.Println("TODO display: ", c.Args().First())
					return nil
				},
			},
			{
				Name:  "version",
				Usage: "version of trams app.",
				Action: func(c *cli.Context) error {
					fmt.Println("version: ", version)
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
