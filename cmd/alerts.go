package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/alerts"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
	scalingo "github.com/Scalingo/go-scalingo/v6"
)

var (
	alertsListCommand = cli.Command{
		Name:      "alerts",
		Category:  "Alerts",
		Flags:     []cli.Flag{&appFlag},
		Usage:     "List the alerts of an application",
		ArgsUsage: " ",
		Description: CommandDescription{
			Description: "List all the alerts of an application",
			Examples:    []string{"scalingo --app my-app alerts"},
			SeeAlso:     []string{"alerts-add", "alerts-update", "alerts-remove"},
		}.Render(),

		Action: func(c *cli.Context) error {
			if c.Args().Len() != 0 {
				cli.ShowCommandHelp(c, "alerts")
				return nil
			}

			err := alerts.List(c.Context, detect.CurrentApp(c))
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "alerts")
		},
	}

	alertsAddCommand = cli.Command{
		Name:     "alerts-add",
		Category: "Alerts",
		Usage:    "Add an alert to an application metric",
		Flags: []cli.Flag{&appFlag,
			&cli.StringFlag{Name: "container-type", Aliases: []string{"c"}, Usage: "Specify the container type affected by the alert", Required: true},
			&cli.StringFlag{Name: "metric", Aliases: []string{"m"}, Usage: "Specify the metric you want the alert to apply on", Required: true},
			&cli.Float64Flag{Name: "limit", Aliases: []string{"l"}, Usage: "Target value for the metric the alert will maintain", Required: true},
			&cli.DurationFlag{Name: "duration-before-trigger", Usage: "Alert is activated if the value is above the limit for the given duration", Required: false},
			&cli.DurationFlag{Name: "remind-every", Aliases: []string{"r"}, Usage: "Send the alert at regular interval when activated", Required: false},
			&cli.BoolFlag{Name: "below", Aliases: []string{"b"}, Usage: "Send the alert when metric value is *below* the limit", Required: false},
			&cli.StringSliceFlag{Name: "notifiers", Aliases: []string{"n"}, Usage: "notifiers' id notified when an alert is activated. Can be specified multiple times.", Required: false},
		},
		Description: CommandDescription{
			Description: "Add an alert to an application metric.\n\nThe 'duration-before-trigger', 'remind-every', 'below' and 'notifiers' flags are optional",
			Examples: []string{
				"scalingo --app my-app alerts-add --container-type web --metric cpu --limit 0.75",
				"scalingo --app my-app alerts-add --container-type web --metric rpm_per_container --limit 100 --remind-every 5m30s --below",
				"scalingo --app my-app alerts-add --container-type web --metric cpu --limit 0.75 --notifiers 5aaab14dcbf5e7000120fd01 --notifiers 5aaab3cacbf5e7000120fd19",
			},
			SeeAlso: []string{"alerts-update", "alerts-remove"},
		}.Render(),
		Action: func(c *cli.Context) error {
			if !isValidAlertAddOpts(c) {
				err := cli.ShowCommandHelp(c, "alerts-add")
				if err != nil {
					errorQuit(err)
				}
				return nil
			}

			currentApp := detect.CurrentApp(c)
			remindEvery := c.Duration("r")
			durationBeforeTrigger := c.Duration("duration-before-trigger")
			err := alerts.Add(c.Context, currentApp, scalingo.AlertAddParams{
				ContainerType:         c.String("c"),
				Metric:                c.String("m"),
				Limit:                 c.Float64("l"),
				SendWhenBelow:         c.Bool("b"),
				DurationBeforeTrigger: &durationBeforeTrigger,
				RemindEvery:           &remindEvery,
				Notifiers:             c.StringSlice("n"),
			})
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "alerts-add")
		},
	}

	alertsUpdateCommand = cli.Command{
		Name:      "alerts-update",
		Category:  "Alerts",
		Usage:     "Update an alert",
		ArgsUsage: "alert-id",
		Flags: []cli.Flag{&appFlag,
			&cli.StringFlag{Name: "container-type", Aliases: []string{"c"}, Usage: "Specify the container type affected by the alert"},
			&cli.StringFlag{Name: "metric", Aliases: []string{"m"}, Usage: "Specify the metric you want the alert to apply on"},
			&cli.Float64Flag{Name: "limit", Aliases: []string{"l"}, Usage: "Target value for the metric the alert will maintain"},
			&cli.DurationFlag{Name: "duration-before-trigger", Usage: "Alert is activated if the value is above the limit for the given duration"},
			&cli.DurationFlag{Name: "remind-every", Aliases: []string{"r"}, Usage: "Send the alert at regular interval when activated"},
			&cli.BoolFlag{Name: "below", Aliases: []string{"b"}, Usage: "Send the alert when metric value is *below* the limit"},
			&cli.BoolFlag{Name: "disabled", Aliases: []string{"d"}, Usage: "Disable the alert (nothing is sent)"},
			&cli.StringSliceFlag{Name: "notifiers", Aliases: []string{"n"}, Usage: "notifiers' id notified when an alert is activated. Can be specified multiple times."},
		},
		Description: CommandDescription{
			Description: "Update an existing alert.\n\nAll flags are optional.",
			Examples: []string{
				"scalingo --app my-app alerts-update --metric rpm-per-container --limit 150 alert-id",
				"scalingo --app my-app alerts-update --disabled alert-id",
				"scalingo --app my-app alerts-update --notifiers 5aaab14dcbf5e7000120fd01 --notifiers 5aaab3cacbf5e7000120fd19 alert-id",
			},
			SeeAlso: []string{"alerts-disable", "alerts-enable"},
		}.Render(),

		Action: func(c *cli.Context) error {
			if c.Args().Len() != 1 {
				err := cli.ShowCommandHelp(c, "alerts-update")
				if err != nil {
					errorQuit(err)
				}
				return nil
			}

			alertID := c.Args().First()
			currentApp := detect.CurrentApp(c)
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
			if c.IsSet("duration-before-trigger") {
				durationBeforeTrigger := c.Duration("duration-before-trigger")
				params.DurationBeforeTrigger = &durationBeforeTrigger
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
			if c.IsSet("n") {
				n := c.StringSlice("n")
				params.Notifiers = &n
			}

			err := alerts.Update(c.Context, currentApp, alertID, params)
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "alerts-add")
		},
	}

	alertsEnableCommand = cli.Command{
		Name:      "alerts-enable",
		Category:  "Alerts",
		Usage:     "Enable an alert",
		ArgsUsage: "alert-id",
		Flags:     []cli.Flag{&appFlag},
		Description: CommandDescription{
			Description: "Enable an alert",
			Examples:    []string{"scalingo --app my-app alerts-enable alert-id"},
			SeeAlso:     []string{"alerts-update", "alerts-remove"},
		}.Render(),

		Action: func(c *cli.Context) error {
			if c.Args().Len() != 1 {
				err := cli.ShowCommandHelp(c, "alerts-enable")
				if err != nil {
					errorQuit(err)
				}
				return nil
			}

			currentApp := detect.CurrentApp(c)
			disabled := false
			err := alerts.Update(c.Context, currentApp, c.Args().First(), scalingo.AlertUpdateParams{
				Disabled: &disabled,
			})
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "alerts-enable")
		},
	}

	alertsDisableCommand = cli.Command{
		Name:      "alerts-disable",
		Category:  "Alerts",
		Usage:     "Disable an alert",
		ArgsUsage: "alert-id",
		Flags:     []cli.Flag{&appFlag},
		Description: CommandDescription{
			Description: "Disable an alert",
			Examples:    []string{"scalingo --app my-app alerts-disable alert-id"},
			SeeAlso:     []string{"alerts-update", "alerts-remove"},
		}.Render(),

		Action: func(c *cli.Context) error {
			if c.Args().Len() != 1 {
				err := cli.ShowCommandHelp(c, "alerts-disable")
				if err != nil {
					errorQuit(err)
				}
				return nil
			}

			currentApp := detect.CurrentApp(c)
			disabled := true
			err := alerts.Update(c.Context, currentApp, c.Args().First(), scalingo.AlertUpdateParams{
				Disabled: &disabled,
			})
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "alerts-disable")
		},
	}

	alertsRemoveCommand = cli.Command{
		Name:      "alerts-remove",
		Category:  "Alerts",
		Usage:     "Remove an alert from an application",
		ArgsUsage: "alert-id",
		Flags:     []cli.Flag{&appFlag},
		Description: CommandDescription{
			Description: "Remove an alert",
			Examples:    []string{"scalingo --app my-app alerts-remove alert-id"},
			SeeAlso:     []string{"alerts-add", "alerts-update"},
		}.Render(),
		Action: func(c *cli.Context) error {
			if c.Args().Len() != 1 {
				cli.ShowCommandHelp(c, "alerts-remove")
				return nil
			}

			currentApp := detect.CurrentApp(c)
			err := alerts.Remove(c.Context, currentApp, c.Args().First())
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "alerts-remove")
		},
	}
)

func isValidAlertAddOpts(c *cli.Context) bool {
	if c.Args().Len() > 0 {
		return false
	}
	for _, opt := range []string{"container-type", "metric", "limit"} {
		if !c.IsSet(opt) {
			return false
		}
	}
	return true
}
