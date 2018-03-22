package cmd

import (
	"github.com/Scalingo/cli/alerts"
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/urfave/cli"
)

var (
	alertsListCommand = cli.Command{
		Name:     "alerts",
		Category: "Alerts",
		Flags:    []cli.Flag{appFlag},
		Usage:    "List the alerts of an application",
		Description: `List all the alerts of an application:

    $ scalingo -a my-app alerts

    # See also commands 'alerts-add' and 'alerts-remove'`,

		Before: AuthenticateHook,
		Action: func(c *cli.Context) {
			if len(c.Args()) != 0 {
				cli.ShowCommandHelp(c, "alerts")
				return
			}

			err := alerts.List(appdetect.CurrentApp(c))
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "alerts")
		},
	}
)
