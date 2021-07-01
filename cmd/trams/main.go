package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"sort"

	"github.com/ayubmalik/trams"
	"github.com/ayubmalik/trams/style"
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
					err := displayMetrolinks(client, c.Args().Slice())
					return err
				},
			},
			{
				Name:  "list",
				Usage: "list all stations with ID and location so can be used by 'display' command.",
				Action: func(c *cli.Context) error {
					err := listStations(client)
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

func displayMetrolinks(client trams.Client, ids []string) error {
	metrolinks, err := client.Get(ids...)
	if err != nil {
		return err
	}

	for _, s := range metrolinks {
		fmt.Printf("%03d %s %s (%s %s)\n", s.Id, s.TLAREF, s.StationLocation, s.PIDREF, s.Direction)
	}
	return nil
}

func listStations(client trams.Client) error {
	stationIDs, err := cachedStations(client)
	if err != nil {
		return err
	}

	uniqueNames := make([]string, 0)
	for _, s := range stationIDs {
		name := fmt.Sprintf("%s %-24s", s.TLAREF, s.StationLocation)
		if !contains(uniqueNames, name) {
			uniqueNames = append(uniqueNames, name)
		}
	}

	sort.Strings(uniqueNames)
	for _, s := range uniqueNames {
		fmt.Println(style.StationName(s))
	}
	return nil
}

func cachedStations(client trams.Client) ([]trams.StationID, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	cache := path.Join(home, ".trams-cache")
	if _, err := os.Stat(cache); os.IsNotExist(err) {
		fmt.Println("TODO")
	}

	return client.List()
}

func contains(values []string, s string) bool {
	for _, v := range values {
		if v == s {
			return true
		}
	}
	return false
}
