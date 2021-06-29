package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ayubmalik/trams"
	"github.com/urfave/cli/v2"
)

const apiURL = "https://europe-west2-tramsfunc.cloudfunctions.net/tramsfunc"

var version = "dev"

func main() {
	client := trams.NewClient(apiURL, 1000)

	app := &cli.App{
		Usage:     "display Metrolink tram information using data from TfGM API",
		UsageText: "trams COMMAND [args]",
		Commands: []*cli.Command{
			{
				Name:      "display",
				Usage:     "display tram information for specified Metrolink stations. If no stations are specified displays all stations. Run `trams help display` for more details.",
				UsageText: "display [STATION...] - If no STATION arguments are specified, displays all stations. Multiple STATION arguments can be specified as short name or long name, e.g. `display BCH MAN` or `display Benchill MAN`",
				Action: func(c *cli.Context) error {
					err := displayStations(client, c.Args().Slice())
					return err
				},
			},
			{
				Name:  "version",
				Usage: "version of trams app.",
				Action: func(c *cli.Context) error {
					fmt.Println(version)
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

func displayStations(client trams.Client, ids []string) error {
	stations, err := client.List(ids...)
	if err != nil {
		return err
	}

	for _, s := range stations {
		fmt.Printf("%03d %s %s (%s %s)\n", s.Id, s.TLAREF, s.StationLocation, s.PIDREF, s.Direction)
	}
	return nil
}
