package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Usage: "display Metrolink tram information using data from TfGM API",
		Action: func(c *cli.Context) error {
			fmt.Println("display stations")
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
