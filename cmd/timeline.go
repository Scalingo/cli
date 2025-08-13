package cmd

import (
	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v8"
)

var (
	TimelineCommand = cli.Command{
		Name:     "timeline",
		Category: "Events",
		Flags: []cli.Flag{
			&appFlag,
			&cli.IntFlag{Name: "page", Usage: "Page to display", Value: 1},
			&cli.IntFlag{Name: "per-page", Usage: "Number of events to display", Value: 30},
		},
		Usage: "List the actions related to a given app",
		Description: CommandDescription{
			Description: "List the actions done by the owner and collaborators of an app",
			Examples:    []string{"scalingo --app my-app timeline"},
		}.Render(),

		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() != 0 {
				cli.ShowCommandHelp(c, "timeline")
				return nil
			}

			utils.CheckForConsent(c.Context, currentApp)
			err := apps.Events(c.Context, currentApp, scalingo.PaginationOpts{
				Page:    c.Int("page"),
				PerPage: c.Int("per-page"),
			})
			if err != nil {
				errorQuit(c.Context, err)
			}

			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "timeline")
		},
	}
)
