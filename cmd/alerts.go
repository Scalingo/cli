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

    # See also commands 'alerts-add', 'alerts-update' and 'alerts-remove'`,

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
			cli.DurationFlag{Name: "remind-every, r", Usage: "When the alert is activated, send the alert at regular interval"},
			cli.BoolFlag{Name: "below, b", Usage: "Send the alert when metric value is *below* the limit"},
		},
		Description: `Add an alert to an application metric.

   The "remind-every" and "below" flags are optionnal

   Example
     scalingo --app my-app alerts-add --container-type web --metric cpu --limit 0.75
     scalingo --app my-app alerts-add --container-type web --metric rpm_per_container --limit 100 --remind-every 5m30s --below

    # See also commands 'alerts-update' and 'alerts-remove'
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
			remindEvery := c.Duration("r")
			err := alerts.Add(currentApp, scalingo.AlertAddParams{
				ContainerType: c.String("c"),
				Metric:        c.String("m"),
				Limit:         c.Float64("l"),
				SendWhenBelow: c.Bool("b"),
				RemindEvery:   &remindEvery,
			})
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "alerts-add")
		},
	}

	alertsUpdateCommand = cli.Command{
		Name:     "alerts-update",
		Category: "Alerts",
		Usage:    "Update an alert",
		Flags: []cli.Flag{appFlag,
			cli.StringFlag{Name: "container-type, c", Usage: "Specify the container type affected by the alert"},
			cli.StringFlag{Name: "metric, m", Usage: "Specify the metric you want the alert to apply on"},
			cli.Float64Flag{Name: "limit, l", Usage: "Target value for the metric the alert will maintain"},
			cli.DurationFlag{Name: "remind-every, r", Usage: "When the alert is activated, send the alert at regular interval"},
			cli.BoolFlag{Name: "below, b", Usage: "Send the alert when metric value is *below* the limit"},
			cli.BoolFlag{Name: "disabled, d", Usage: "Disable the alert (nothing is sent)"},
		},
		Description: `Update an existing alert.

   All flags are optionnal.

   Example
     scalingo --app my-app alerts-update --metric rpm-per-container --target 150 <ID>
     scalingo --app my-app alerts-update --disabled <ID>

   # See also 'alerts-disable' and 'alerts-enable'
		`,
		Before: AuthenticateHook,
		Action: func(c *cli.Context) {
			if len(c.Args()) != 1 {
				err := cli.ShowCommandHelp(c, "alerts-update")
				if err != nil {
					errorQuit(err)
				}
				return
			}

			alertID := c.Args()[0]
			currentApp := appdetect.CurrentApp(c)
			params := scalingo.AlertUpdateParams{}
			if c.IsSet("c") {
				ct := c.String("c")
				params.ContainerType = &ct
			}
			if c.IsSet("m") {
				m := c.String("m")
				params.Metric = &m
			}
			if c.IsSet("l") {
				l := c.Float64("l")
				params.Limit = &l
			}
			if c.IsSet("r") {
				remindEvery := c.Duration("r")
				params.RemindEvery = &remindEvery
			}
			if c.IsSet("b") {
				b := c.Bool("b")
				params.SendWhenBelow = &b
			}
			if c.IsSet("d") {
				d := c.Bool("d")
				params.Disabled = &d
			}

			err := alerts.Update(currentApp, alertID, params)
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "alerts-add")
		},
	}

	alertsEnableCommand = cli.Command{
		Name:     "alerts-enable",
		Category: "Alerts",
		Usage:    "Enable an alert",
		Flags:    []cli.Flag{appFlag},
		Description: `Enable an alert.

   Example
     scalingo --app my-app alerts-enable <ID>

   # See also commands 'alerts-update' and 'alerts-remove'
		`,
		Before: AuthenticateHook,
		Action: func(c *cli.Context) {
			if len(c.Args()) != 1 {
				err := cli.ShowCommandHelp(c, "alerts-enable")
				if err != nil {
					errorQuit(err)
				}
				return
			}

			currentApp := appdetect.CurrentApp(c)
			disabled := false
			err := alerts.Update(currentApp, c.Args()[0], scalingo.AlertUpdateParams{
				Disabled: &disabled,
			})
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "alerts-enable")
		},
	}

	alertsDisableCommand = cli.Command{
		Name:     "alerts-disable",
		Category: "Alerts",
		Usage:    "Disable an alert",
		Flags:    []cli.Flag{appFlag},
		Description: `Disable an alert.

   Example
     scalingo --app my-app alerts-disable <ID>

   # See also commands 'alerts-update' and 'alerts-remove'
		`,
		Before: AuthenticateHook,
		Action: func(c *cli.Context) {
			if len(c.Args()) != 1 {
				err := cli.ShowCommandHelp(c, "alerts-disable")
				if err != nil {
					errorQuit(err)
				}
				return
			}

			currentApp := appdetect.CurrentApp(c)
			disabled := true
			err := alerts.Update(currentApp, c.Args()[0], scalingo.AlertUpdateParams{
				Disabled: &disabled,
			})
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "alerts-disable")
		},
	}

	alertsRemoveCommand = cli.Command{
		Name:     "alerts-remove",
		Category: "Alerts",
		Usage:    "Remove an alert from an application",
		Flags:    []cli.Flag{appFlag},
		Description: `Remove an alert.

   Example
     scalingo --app my-app alerts-remove <ID>

   # See also commands 'alerts-add' and 'alerts-update'
		 `,
		Before: AuthenticateHook,
		Action: func(c *cli.Context) {
			if len(c.Args()) != 1 {
				cli.ShowCommandHelp(c, "alerts-remove")
				return
			}

			currentApp := appdetect.CurrentApp(c)
			err := alerts.Remove(currentApp, c.Args()[0])
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "alerts-remove")
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
