package cmd

import (
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/crontasks"
	"github.com/urfave/cli"
)

var (
	cronTasksListCommand = cli.Command{
		Name:     "cron-jobs",
		Category: "Cron Jobs",
		Flags:    []cli.Flag{appFlag},
		Usage:    "List the cron jobs of an application",
		Description: `List all the cron jobs of an application:

    $ scalingo --app my-app cron-jobs`,

		Action: func(c *cli.Context) {
			if len(c.Args()) > 0 {
				cli.ShowCommandHelp(c, "cron-jobs")
				return
			}

			currentApp := appdetect.CurrentApp(c)
			err := crontasks.List(currentApp)
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "cron-jobs")
		},
	}
)
