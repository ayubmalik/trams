package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"strconv"

	"github.com/ayubmalik/trams"
	"github.com/ayubmalik/trams/style"
	"github.com/urfave/cli/v2"
)

const (
	apiURL = "https://europe-west2-tramsfunc.cloudfunctions.net/tramsfunc"
)

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

func displayMetrolinks(client trams.Client, refs []string) error {
	ids := lookupIDs(client, refs)
	metrolinks, err := client.Get(ids...)
	if err != nil {
		return err
	}

	for _, m := range metrolinks {
		fmt.Println(style.FormatMetrolink(m, 0))
	}
	return nil
}

func lookupIDs(client trams.Client, refs []string) []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	cache := path.Join(home, ".trams-cache")

	grouped, err := cachedStations(client, cache)
	if err != nil {
		return nil
	}

	ids := make([]string, 0)
	for _, r := range refs {
		stationIDs := grouped[r]
		for _, s := range stationIDs {
			ids = append(ids, strconv.Itoa(s.Id))
		}
	}
	return ids
}

func listStations(client trams.Client) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	cache := path.Join(home, ".trams-cache")

	stationIDs, err := cachedStations(client, cache)
	if err != nil {
		return err
	}

	uniqueNames := make([]string, 0)
	for _, s := range stationIDs {
		// TODO: move to style pkg
		name := fmt.Sprintf("| %s | %-40s", s[0].TLAREF, s[0].StationLocation)
		uniqueNames = append(uniqueNames, name)
	}

	sort.Strings(uniqueNames)
	rows := style.StationIDRows(uniqueNames)
	for _, r := range rows {
		fmt.Println(r)
	}
	return nil
}

func cachedStations(client trams.Client, cache string) (map[string][]trams.StationID, error) {
	var stationIDs []trams.StationID
	var grouped map[string][]trams.StationID

	if _, err := os.Stat(cache); os.IsNotExist(err) {
		stationIDs, err = client.List()
		if err != nil {
			return nil, err
		}

		f, err := os.Create(cache)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		json.NewEncoder(f).Encode(stationIDs)
	} else {
		f, err := os.Open(cache)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		json.NewDecoder(f).Decode(&stationIDs)
	}

	grouped = groupStationsByRef(stationIDs)
	return grouped, nil
}

func groupStationsByRef(stationIDS []trams.StationID) map[string][]trams.StationID {
	gs := make(map[string][]trams.StationID)
	for _, s := range stationIDS {
		gs[s.TLAREF] = append(gs[s.TLAREF], s)
	}
	return gs
}
