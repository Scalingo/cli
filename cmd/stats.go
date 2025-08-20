package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

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

		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() != 0 {
				_ = cli.ShowCommandHelp(ctx, c, "stats")
				return nil
			}

			err := apps.Stats(ctx, currentApp, c.Bool("stream"))
			if err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
	}
)
