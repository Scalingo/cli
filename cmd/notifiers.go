package cmd

import (
	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/notifiers"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v8"
)

var (
	NotifiersListCommand = cli.Command{
		Name:     "notifiers",
		Category: "Notifiers",
		Usage:    "List your notifiers",
		Flags:    []cli.Flag{&appFlag},
		Description: CommandDescription{
			Description: "List all notifiers of your app",
			Examples:    []string{"scalingo --app my-app notifiers"},
			SeeAlso:     []string{"notifiers-add", "notifiers-remove"},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)
			var err error
			if c.Args().Len() == 0 {
				err = notifiers.List(c.Context, currentApp)
			} else {
				cli.ShowCommandHelp(c, "notifiers")
			}

			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "notifiers")
		},
	}

	NotifiersDetailsCommand = cli.Command{
		Name:      "notifiers-details",
		Category:  "Notifiers",
		Usage:     "Show details of a notifier",
		ArgsUsage: "notifier-id",
		Flags:     []cli.Flag{&appFlag},
		Description: CommandDescription{
			Description: "Show details of a notifier",
			Examples:    []string{"scalingo --app my-app notifiers-details my-notifier"},
			SeeAlso:     []string{"notifiers"},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)
			var err error
			if c.Args().Len() == 1 {
				err = notifiers.Details(c.Context, currentApp, c.Args().First())
			} else {
				cli.ShowCommandHelp(c, "notifiers-details")
			}

			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "notifiers-details")
			autocomplete.NotifiersAutoComplete(c)
		},
	}

	NotifiersAddCommand = cli.Command{
		Name:     "notifiers-add",
		Category: "Notifiers",
		Flags: []cli.Flag{
			&appFlag,
			&cli.BoolFlag{Name: "enable", Aliases: []string{"e"}, Usage: "Enable the notifier (default)"},
			&cli.BoolFlag{Name: "disable", Aliases: []string{"d"}, Usage: "Disable the notifier"},
			&cli.StringFlag{Name: "platform", Aliases: []string{"p"}, Value: "", Usage: "The notifier platform"},
			&cli.StringFlag{Name: "name", Aliases: []string{"n"}, Value: "", Usage: "Name of the notifier"},
			&cli.BoolFlag{Name: "send-all-events", Aliases: []string{"sa"}, Usage: "If true the notifier will send all events. Default: false"},
			&cli.StringFlag{Name: "webhook-url", Aliases: []string{"u"}, Value: "", Usage: "The webhook url to send notification (if applicable)"},
			&cli.StringFlag{Name: "phone", Value: "", Usage: "The phone number to send notifications (if applicable)"},
			&cli.StringSliceFlag{Name: "event", Aliases: []string{"ev"}, Value: cli.NewStringSlice(), Usage: "List of selected events. Default: []"},
			&cli.StringSliceFlag{Name: "email", Value: cli.NewStringSlice(), Usage: "The emails (multiple option accepted) to send notifications (if applicable)"},
			&cli.StringSliceFlag{Name: "collaborator", Value: cli.NewStringSlice(), Usage: "The usernames of the collaborators who will receive notifications"},
		},
		Usage: "Add a notifier for your application",
		Description: CommandDescription{
			Description: "Add a notifier for your application",
			Examples: []string{
				"scalingo --app my-app notifiers-add --platform slack --name \"My notifier\" --webhook-url https://hooks.slack.com/services/1234 --event deployment --event stop_app",
				"scalingo --app my-app notifiers-add --platform webhook --name \"My notifier\" --webhook-url https://custom-webhook.com --send-all-events",
			},
			SeeAlso: []string{"notifiers", "notifiers-remove"},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)

			utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeContainers)

			if c.String("platform") == "" {
				cli.ShowCommandHelp(c, "notifiers-add")
			}

			var active bool
			if c.IsSet("disable") {
				active = false
			} else {
				active = true
			}
			sendAllEvents := c.Bool("send-all-events")

			params := notifiers.ProvisionParams{
				CollaboratorUsernames: c.StringSlice("collaborator"),
				SelectedEventNames:    c.StringSlice("event"),
				NotifierParams: scalingo.NotifierParams{
					Active:        &active,
					Name:          c.String("name"),
					SendAllEvents: &sendAllEvents,

					// Type data options
					PhoneNumber: c.String("phone"),
					Emails:      c.StringSlice("email"),
					WebhookURL:  c.String("webhook-url"),
				},
			}

			err := notifiers.Provision(c.Context, currentApp, c.String("platform"), params)
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "notifiers-add")
		},
	}

	NotifiersUpdateCommand = cli.Command{
		Name:     "notifiers-update",
		Category: "Notifiers",
		Flags: []cli.Flag{
			&appFlag,
			&cli.BoolFlag{Name: "enable", Aliases: []string{"e"}, Usage: "Enable the notifier"},
			&cli.BoolFlag{Name: "disable", Aliases: []string{"d"}, Usage: "Disable the notifier"},
			&cli.StringFlag{Name: "name", Aliases: []string{"n"}, Value: "", Usage: "Name of the notifier"},
			&cli.BoolFlag{Name: "send-all-events", Aliases: []string{"sa"}, Usage: "If true the notifier will send all events. Default: false"},
			&cli.StringFlag{Name: "webhook-url", Aliases: []string{"u"}, Value: "", Usage: "The webhook url to send notification (if applicable)"},
			&cli.StringFlag{Name: "phone", Value: "", Usage: "The phone number to send notifications (if applicable)"},
			&cli.StringFlag{Name: "email", Value: "", Usage: "The email to send notifications (if applicable)"},
			&cli.StringSliceFlag{Name: "event", Aliases: []string{"ev"}, Value: cli.NewStringSlice(), Usage: "List of selected events. Default: []"},
		},
		Usage:     "Update a notifier",
		ArgsUsage: "notifier-id",
		Description: CommandDescription{
			Description: "Update a notifier",
			Examples: []string{
				"scalingo -a myapp notifiers-update --disable my-notifier",
				"scalingo -a myapp notifiers-update --name \"My notifier\" --webhook-url https://custom-webhook.com --send-all-events my-notifier",
			},
			SeeAlso: []string{"notifiers", "notifiers-remove"},
		}.Render(),

		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)
			utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeContainers)
			var err error

			var active *bool
			if c.IsSet("enable") {
				active = utils.BoolPtr(true)
			} else if c.IsSet("disable") {
				active = utils.BoolPtr(false)
			} else {
				active = nil
			}

			var sendAllEvents *bool
			if c.IsSet("send-all-events") {
				tmpEvents := c.Bool("send-all-events")
				sendAllEvents = &tmpEvents
			} else {
				sendAllEvents = nil
			}

			params := notifiers.ProvisionParams{
				CollaboratorUsernames: c.StringSlice("collaborator"),
				SelectedEventNames:    c.StringSlice("event"),
				NotifierParams: scalingo.NotifierParams{
					Active:        active,
					Name:          c.String("name"),
					SendAllEvents: sendAllEvents,

					// Type data options
					PhoneNumber: c.String("phone"),
					Emails:      c.StringSlice("email"),
					WebhookURL:  c.String("webhook-url"),
				},
			}
			if c.Args().Len() >= 1 {
				err = notifiers.Update(c.Context, currentApp, c.Args().First(), params)
			} else {
				cli.ShowCommandHelp(c, "notifiers-update")
			}
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "notifiers-update")
			autocomplete.NotifiersAutoComplete(c)
		},
	}

	NotifiersRemoveCommand = cli.Command{
		Name:      "notifiers-remove",
		Category:  "Notifiers",
		Flags:     []cli.Flag{&appFlag},
		Usage:     "Remove an existing notifier from your app",
		ArgsUsage: "notifier-id",
		Description: CommandDescription{
			Description: "Remove an existing notifier from your app",
			Examples:    []string{"scalingo --app my-app notifier-remove my-notifier"},
			SeeAlso:     []string{"notifiers", "notifiers-add"},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)
			utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeContainers)
			var err error
			if c.Args().Len() == 1 {
				err = notifiers.Destroy(c.Context, currentApp, c.Args().First())
			} else {
				cli.ShowCommandHelp(c, "notifiers-remove")
			}
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "notifiers-remove")
			autocomplete.NotifiersAutoComplete(c)
		},
	}
)
