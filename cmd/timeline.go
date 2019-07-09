package cmd

import (
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/go-scalingo"
	"github.com/urfave/cli"
)

var (
	TimelineCommand = cli.Command{
		Name:     "timeline",
		Category: "Events",
		Flags: []cli.Flag{
			appFlag,
			cli.IntFlag{Name: "page", Usage: "Page to display", Value: 1},
			cli.IntFlag{Name: "per-page", Usage: "Number of events to display", Value: 30},
		},
		Usage: "List the actions related to a given app",
		Description: `List the actions done by the owner and collaborators of an app:

    $ scalingo -a myapp timeline`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			var err error
			if len(c.Args()) == 0 {
				err = apps.Events(currentApp, scalingo.PaginationOpts{
					Page:    c.Int("page"),
					PerPage: c.Int("per-page"),
				})
			} else {
				cli.ShowCommandHelp(c, "timeline")
			}

			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "timeline")
		},
	}
)
