package cmd

import (
	"github.com/urfave/cli"

	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/crontasks"
)

var (
	cronTasksListCommand = cli.Command{
		Name:     "cron-tasks",
		Category: "Cron Tasks",
		Flags:    []cli.Flag{appFlag},
		Usage:    "List the cron tasks of an application",
		Description: `List all the cron tasks of an application:

    $ scalingo --app my-app cron-jobs`,

		Action: func(c *cli.Context) {
			if len(c.Args()) > 0 {
				cli.ShowCommandHelp(c, "cron-tasks")
				return
			}

			currentApp := appdetect.CurrentApp(c)
			err := crontasks.List(currentApp)
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "cron-tasks")
		},
	}
)
