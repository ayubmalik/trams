package main

import (
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

	app := &cli.App{
		Usage:     "Metrolink tram information using data from TfGM API. Use the appropriate sub command (see below) e.g. 'trams list' or 'trams display'",
		UsageText: "trams COMMAND [args]",
		Commands: []*cli.Command{
			{
				Name:      "display",
				Usage:     "display tram information for specified Metrolink stations. e.g. 'trams display ABM'. Multiple stations are separated by a space. If no stations are specified displays all stations. Run 'trams help display' for more details.",
				UsageText: "display [STATION...] - If no STATION arguments are specified, displays all stations. Multiple STATION arguments can be specified as short name e.g. 'trams display BCH MAN VIC'",
				Action: func(c *cli.Context) error {
					err := displayMetrolinks(c.Args().Slice())
					return err
				},
			},
			{
				Name:  "list",
				Usage: "list all stations with short reference and name. The short reference is supplied to the 'display' command.",
				Action: func(c *cli.Context) error {
					err := listStations()
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

func displayMetrolinks(refs []string) error {
	groupedStationIDs, err := getAllStationsGroupedByRef()
	if err != nil {
		return err
	}
	ids := lookupIDs(groupedStationIDs, refs)

	client := trams.NewClient(apiURL, 2000)
	metrolinks, err := client.Get(ids...)
	if err != nil {
		return err
	}

	rows := style.MetrolinkRows(metrolinks)
	for _, r := range rows {
		fmt.Println(r)
	}
	return nil
}

func lookupIDs(groupedStationIDs map[string][]trams.StationID, refs []string) []string {
	ids := make([]string, 0)
	for _, r := range refs {
		stationIDs := groupedStationIDs[r]
		for _, s := range stationIDs {
			ids = append(ids, strconv.Itoa(s.Id))
		}
	}
	return ids
}

func listStations() error {
	groupedStationIDs, err := getAllStationsGroupedByRef()
	if err != nil {
		return err
	}

	uniqueNames := make([]string, 0)
	for _, s := range groupedStationIDs {
		name := style.FormatStationID(s[0])
		uniqueNames = append(uniqueNames, name)
	}

	sort.Strings(uniqueNames)
	rows := style.StationIDRows(uniqueNames)
	for _, r := range rows {
		fmt.Println(r)
	}
	return nil
}

func getCacheFile() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	cache := path.Join(home, ".trams-cache")
	return cache, nil
}

func getAllStationsGroupedByRef() (map[string][]trams.StationID, error) {
	client := trams.NewClient(apiURL, 1000)
	stationIDs, err := client.List()
	if err != nil {
		return nil, err
	}

	groupedStationIDs := groupStationIDsByRef(stationIDs)
	return groupedStationIDs, nil
}

func groupStationIDsByRef(stationIDS []trams.StationID) map[string][]trams.StationID {
	gs := make(map[string][]trams.StationID)
	for _, s := range stationIDS {
		gs[s.TLAREF] = append(gs[s.TLAREF], s)
	}
	return gs
}
