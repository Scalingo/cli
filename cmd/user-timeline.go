package cmd

import (
	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/user"
	"github.com/Scalingo/go-scalingo/v8"
)

var (
	UserTimelineCommand = cli.Command{
		Name:     "user-timeline",
		Category: "Events",
		Flags: []cli.Flag{
			&cli.IntFlag{Name: "page", Usage: "Page to display", Value: 1},
			&cli.IntFlag{Name: "per-page", Usage: "Number of events to display", Value: 30},
		},
		Usage: "List the events you have done on the platform",
		Description: CommandDescription{
			Description: "List the events you have done on the platform",
			Examples:    []string{"scalingo user-timeline --page 3 --per-page 20"},
		}.Render(),

		Action: func(ctx context.Context, c *cli.Command) error {
			var err error
			if c.Args().Len() == 0 {
				err = user.Events(c.Context, scalingo.PaginationOpts{
					Page:    c.Int("page"),
					PerPage: c.Int("per-page"),
				})
			} else {
				cli.ShowCommandHelp(c, "user-timeline")
			}

			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "user-timeline")
		},
	}
)
