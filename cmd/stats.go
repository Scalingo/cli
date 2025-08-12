package cmd

import (
	cli "github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/detect"
)

var (
	StatsCommand = cli.Command{
		Name:     "stats",
		Category: "Display metrics of the running containers",
		Usage:    "Display metrics of the currently running containers",
		Flags: []cli.Flag{
			&appFlag,
			&cli.BoolFlag{Name: "stream", Usage: "Stream metrics data"},
		},
		Description: CommandDescription{
			Description: "Display metrics of your application running containers",
			Examples:    []string{"scalingo --app my-app stats"},
		}.Render(),

		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() != 0 {
				cli.ShowCommandHelp(c, "stats")
				return nil
			}

			err := apps.Stats(c.Context, currentApp, c.Bool("stream"))
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {},
	}
)
