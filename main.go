package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v2"
)

func init() {
	log.SetOutput(os.Stderr)
	log.SetLevel(log.InfoLevel)
	log.SetReportTimestamp(false)
	log.SetTimeFormat("2006.01.02 15:04:05")
}

func main() {
	buildInfo, _ := debug.ReadBuildInfo()
	app := &cli.App{
		Name:  "goup",
		Usage: "update dependencies for all projects at once",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "enable debug mode",
				Action: func(_ *cli.Context, debug bool) error {
					if debug {
						log.SetLevel(log.DebugLevel)
						log.SetReportTimestamp(true)
					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"vv"},
				Usage:   "verbose output",
			},
			&cli.BoolFlag{
				Name:    "recursive",
				Aliases: []string{"r"},
				Usage:   "recursively walk directories",
			},
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "include hidden directories",
			},
			&cli.BoolFlag{
				Name:  "vendor",
				Usage: "include vendor directory",
			},
			&cli.StringSliceFlag{
				Name:    "exclude",
				Aliases: []string{"e"},
				Usage:   "exclude directories that match pattern",
			},
		},
		ArgsUsage:            "[dirs]...",
		EnableBashCompletion: true,
		HideHelpCommand:      true,
		BashComplete:         cli.DefaultAppComplete,
		Action:               run,
		Authors: []*cli.Author{
			{
				Name:  "Artem Yadelskyi",
				Email: "mymmrac@gmail.com",
			},
		},
		Suggest: true,
		Version: fmt.Sprintf("%s build with %s", buildInfo.Main.Version, buildInfo.GoVersion),
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
