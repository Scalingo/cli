package cmd

import (
	"github.com/Scalingo/cli/alerts"
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/cmd/autocomplete"
	scalingo "github.com/Scalingo/go-scalingo"
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

	alertsAddCommand = cli.Command{
		Name:     "alerts-add",
		Category: "Alerts",
		Usage:    "Add an alert to an application",
		Flags: []cli.Flag{appFlag,
			cli.StringFlag{Name: "container-type, c", Usage: "Specify the container type affected by the alert"},
			cli.StringFlag{Name: "metric, m", Usage: "Specify the metric you want the alert to apply on"},
			cli.Float64Flag{Name: "limit, l", Usage: "Target value for the metric the alert will maintain"},
			cli.BoolFlag{Name: "disabled, d", Usage: "Disable the alert (nothing is sent)"},
		},
		Description: `Add an alert to an application metric.

   The "disabled" flag is optionnal

   Example
     scalingo --app my-app alerts-add --container-type web --metric cpu --target 0.75 [--disabled]
		`,
		Before: AuthenticateHook,
		Action: func(c *cli.Context) {
			if !isValidAlertAddOpts(c) {
				err := cli.ShowCommandHelp(c, "alerts-add")
				if err != nil {
					errorQuit(err)
				}
				return
			}

			currentApp := appdetect.CurrentApp(c)
			err := alerts.Add(currentApp, scalingo.AlertParams{
				ContainerType: c.String("c"),
				Metric:        c.String("m"),
				Limit:         c.Float64("l"),
				Disabled:      c.Bool("d"),
			})
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "alerts-add")
		},
	}
)

func isValidAlertAddOpts(c *cli.Context) bool {
	if len(c.Args()) > 0 {
		return false
	}
	for _, opt := range []string{"container-type", "metric", "limit"} {
		if !c.IsSet(opt) {
			return false
		}
	}
	return true
}
