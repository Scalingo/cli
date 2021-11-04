package cmd

import (
	"github.com/urfave/cli"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/user"
	"github.com/Scalingo/go-scalingo/v4"
)

var (
	UserTimelineCommand = cli.Command{
		Name:     "user-timeline",
		Category: "Events",
		Flags: []cli.Flag{
			cli.IntFlag{Name: "page", Usage: "Page to display", Value: 1},
			cli.IntFlag{Name: "per-page", Usage: "Number of events to display", Value: 30},
		},
		Usage: "List the events you have done on the platform",
		Description: `List the events you have done on the platform:

    $ scalingo user-timeline`,
		Action: func(c *cli.Context) {
			var err error
			if len(c.Args()) == 0 {
				err = user.Events(scalingo.PaginationOpts{
					Page:    c.Int("page"),
					PerPage: c.Int("per-page"),
				})
			} else {
				cli.ShowCommandHelp(c, "user-timeline")
			}

			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "user-timeline")
		},
	}
)
