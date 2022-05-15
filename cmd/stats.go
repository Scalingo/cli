package cmd

import (
	"github.com/urfave/cli"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/detect"
)

var (
	StatsCommand = cli.Command{
		Name:     "stats",
		Category: "Display metrics of the running containers",
		Usage:    "Display metrics of the currently running containers",
		Flags: []cli.Flag{
			appFlag,
			cli.BoolFlag{Name: "stream", Usage: "Stream metrics data"},
		},
		Description: `Display metrics of you application running containers
	Example
	  'scalingo --app my-app stats'`,
		Action: func(c *cli.Context) {
			currentApp := detect.CurrentApp(c)
			if len(c.Args()) != 0 {
				cli.ShowCommandHelp(c, "stats")
			} else if err := apps.Stats(currentApp, c.Bool("stream")); err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			return
		},
	}
)
