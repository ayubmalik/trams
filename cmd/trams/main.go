package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ayubmalik/trams"
	"github.com/charmbracelet/lipgloss"
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
	metrolinks, err := client.List(ids...)
	if err != nil {
		return err
	}

	for _, m := range metrolinks {
		fmt.Println(FormatMetrolink(m))
	}
	return nil
}

func FormatMetrolink(m trams.Metrolink) string {
	pad := 1
	width := 52

	style := lipgloss.NewStyle().
		Bold(false).
		PaddingLeft(pad).
		PaddingRight(pad).
		Foreground(lipgloss.Color("215")).
		Background(lipgloss.Color("234")).
		Width(width).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("199"))

	inline := lipgloss.NewStyle().
		Inline(true).
		Bold(true).
		Foreground(lipgloss.Color("16")).
		Background(lipgloss.Color("215")).
		Width(width)

	text := fmt.Sprintf("%s %s (Platform %s %s)", m.TLAREF, m.StationLocation, m.Platform(), m.Direction)
	text = inline.Render(text)

	if m.Status0 == "" {
		text += "\nNo information available"
	}
	if m.Status0 != "" {
		text += fmt.Sprintf("\n%sm %s", m.Wait0, m.Dest0)
	}
	if m.Status1 != "" {
		text += fmt.Sprintf("\n%sm %s", m.Wait1, m.Dest1)
	}

	if m.Status2 != "" {
		text += fmt.Sprintf("\n%sm %s", m.Wait2, m.Dest2)
	}

	inline2 := lipgloss.NewStyle().Inline(true).Foreground(lipgloss.Color("230"))
	text += "\n\n" + inline2.Render(m.MessageBoard)
	return style.Render(text)

}
