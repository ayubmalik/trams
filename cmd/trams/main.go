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

	app := &cli.App{
		Usage:     "display Metrolink tram information using data from TfGM API",
		UsageText: "trams COMMAND [args]",
		Commands: []*cli.Command{
			{
				Name:      "display",
				Usage:     "display tram information for specified Metrolink stations. If no stations are specified displays all stations. Run `trams help display` for more details.",
				UsageText: "display [STATION...] - If no STATION arguments are specified, displays all stations. Multiple STATION arguments can be specified as short name e.g. `display BCH MAN VIC`",
				Action: func(c *cli.Context) error {
					err := displayMetrolinks(c.Args().Slice())
					return err
				},
			},
			{
				Name:  "list",
				Usage: "list all stations with short name (TLAREF) and name so can be used by 'display' command.",
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

	client := trams.NewClient(apiURL, 1000)
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
	cache, err := getCacheFile()
	if err != nil {
		return nil, err
	}

	groupedStationIDs, err := cachedStations(client, cache)
	if err != nil {
		return nil, err
	}
	return groupedStationIDs, nil
}

func cachedStations(client trams.Client, cache string) (map[string][]trams.StationID, error) {
	var stationIDs []trams.StationID
	var groupedStationIDs map[string][]trams.StationID

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
	}

	f, err := os.Open(cache)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	json.NewDecoder(f).Decode(&stationIDs)

	groupedStationIDs = groupStationIDsByRef(stationIDs)
	return groupedStationIDs, nil
}

func groupStationIDsByRef(stationIDS []trams.StationID) map[string][]trams.StationID {
	gs := make(map[string][]trams.StationID)
	for _, s := range stationIDS {
		gs[s.TLAREF] = append(gs[s.TLAREF], s)
	}
	return gs
}
