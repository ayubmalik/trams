package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"time"

	"github.com/ayubmalik/trams"
	"github.com/ayubmalik/trams/style"
	"github.com/rivo/tview"
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
					DisplayUI()
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

func DisplayUI() {
	app := tview.NewApplication()
	flex := tview.NewFlex().
		AddItem(tview.NewBox().SetBorder(true).SetTitle("1 (1/2 x width of Top)"), 0, 1, false).
		AddItem(tview.NewBox().SetBorder(true).SetTitle("2 (1/2 x width of Top)"), 0, 1, false).
		AddItem(tview.NewBox().SetBorder(true).SetTitle("2 (1/2 x width of Top)"), 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Top"), 0, 1, false).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Middle (3 x height of Top)"), 0, 3, false).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Bottom (5 rows)"), 5, 1, false), 0, 2, false).
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Right (20 cols)"), 20, 1, false)

	go func() {
		time.Sleep(50 * time.Millisecond)
		app.EnableMouse(false)
		fmt.Println("\x1b[31mBoom\x1b[0mShoom1")
		fmt.Printf("\x1b[?20h\n")
		fmt.Println("What")
		fmt.Println("The")
		fmt.Println("Funk")
		os.Exit(0)
	}()

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
