package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/crontasks"
	"github.com/Scalingo/cli/detect"
)

var (
	cronTasksListCommand = cli.Command{
		Name:     "cron-tasks",
		Category: "Cron Tasks",
		Flags:    []cli.Flag{&appFlag},
		Usage:    "List the cron tasks of an application",
		Description: CommandDescription{
			Description: "List all the cron tasks of an application",
			Examples:    []string{"scalingo --app my-app cron-tasks"},
		}.Render(),

		Action: func(c *cli.Context) error {
			if c.Args().Len() > 0 {
				cli.ShowCommandHelp(c, "cron-tasks")
				return nil
			}

			currentApp := detect.CurrentApp(c)
			err := crontasks.List(c.Context, currentApp)
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "cron-tasks")
		},
	}
)
