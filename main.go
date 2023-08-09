package main

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v2"
)

func init() {
	log.SetOutput(os.Stderr)
	log.SetLevel(log.InfoLevel)
	log.SetReportTimestamp(true)
	log.SetTimeFormat("2006.01.02 15:04:05")
}

func main() {
	app := &cli.App{
		Name:  "goup",
		Usage: "Update dependencies for all projects at once",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "enable debug mode",
				Action: func(_ *cli.Context, debug bool) error {
					if debug {
						log.SetLevel(log.DebugLevel)
					}
					return nil
				},
			},
		},
		ArgsUsage:            "[dirs...] paths of directories that should be update",
		EnableBashCompletion: true,
		BashComplete:         cli.DefaultAppComplete,
		Action:               run,
		Authors: []*cli.Author{
			{
				Name:  "Artem Yadelskyi",
				Email: "mymmrac@gmail.com",
			},
		},
		Suggest: true,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
