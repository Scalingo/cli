package cmd

import (
	"strings"

	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/notifiers"
	"github.com/Scalingo/codegangsta-cli"
	scalingo "github.com/Scalingo/go-scalingo"
)

var (
	NotifiersListCommand = cli.Command{
		Name:     "notifiers",
		Category: "Notifiers",
		Usage:    "List your notifiers",
		Flags:    []cli.Flag{appFlag},
		Description: ` List all notifiers of your app:
    $ scalingo -a myapp notifiers

		# See also 'notifiers-add' and 'notifiers-remove'
`,
		Before: AuthenticateHook,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			var err error
			if len(c.Args()) == 0 {
				err = notifiers.List(currentApp)
			} else {
				cli.ShowCommandHelp(c, "notifiers")
			}

			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "notifiers")
		},
	}

	NotifiersAddCommand = cli.Command{
		Name:     "notifiers-add",
		Category: "Notifiers",
		Flags: []cli.Flag{
			appFlag,
			cli.BoolFlag{Name: "enable, e", Usage: "Enable the notifier (default)"},
			cli.BoolFlag{Name: "disable, d", Usage: "Disable the notifier"},
			cli.StringFlag{Name: "platform, p", Value: "", Usage: "The notifier platform"},
			cli.StringFlag{Name: "name, n", Value: "", Usage: "Name of the notifier"},
			cli.BoolFlag{Name: "send-all-events, sa", Usage: "If true the notifier will send all events. Default: false"},
			cli.StringFlag{Name: "events, ev", Value: "", Usage: "List of selected events. Default: []"},
			cli.StringFlag{Name: "webhook-url, u", Value: "", Usage: "The webhook url to send notification (if applicable)"},
			cli.StringFlag{Name: "phone", Value: "", Usage: "The phone number to send notifications (if applicable)"},
			cli.StringFlag{Name: "email", Value: "", Usage: "The email to send notifications (if applicable)"},
		},
		Usage: "Add a notifier for your application",
		Description: `Add a notifier for your application:

Examples
 $ scalingo -a myapp notifiers-add \
   --platform slack \
   --name "My notifier" \
   --webhook-url "https://hooks.slack.com/services/1234" \
   --events "deployment stop_app"

 $ scalingo -a myapp notifiers-add \
   --platform webhook \
   --name "My notifier" \
   --webhook-url "https://custom-webhook.com" \
   --send-all-events true

 # See also 'notifiers' and 'notifiers-remove'
`,
		Before: AuthenticateHook,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)

			var active bool
			if c.IsSet("disable") {
				active = false
			} else {
				active = true
			}
			sendAllEvents := c.Bool("send-all-events")

			params := scalingo.NotifierParams{
				Active:         &active,
				Name:           c.String("name"),
				SendAllEvents:  &sendAllEvents,
				SelectedEvents: strings.Split(c.String("events"), " "),

				// Type data options
				PhoneNumber: c.String("phone"),
				Email:       c.String("email"),
				WebhookURL:  c.String("webhook-url"),
			}

			err := notifiers.Provision(currentApp, c.String("platform"), params)
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "notifiers-add")
		},
	}

	NotifiersUpdateCommand = cli.Command{
		Name:     "notifiers-update",
		Category: "Notifiers",
		Flags: []cli.Flag{
			appFlag,
			cli.BoolFlag{Name: "enable, e", Usage: "Enable the notifier"},
			cli.BoolFlag{Name: "disable, d", Usage: "Disable the notifier"},
			cli.StringFlag{Name: "name, n", Value: "", Usage: "Name of the notifier"},
			cli.BoolFlag{Name: "send-all-events, sa", Usage: "If true the notifier will send all events. Default: false"},
			cli.StringFlag{Name: "events, ev", Value: "", Usage: "List of selected events. Default: []"},
			cli.StringFlag{Name: "webhook-url, u", Value: "", Usage: "The webhook url to send notification (if applicable)"},
			cli.StringFlag{Name: "phone", Value: "", Usage: "The phone number to send notifications (if applicable)"},
			cli.StringFlag{Name: "email", Value: "", Usage: "The email to send notifications (if applicable)"},
		},
		Usage: "Update a notifier",
		Description: `Update a notifier:
Examples
 $ scalingo -a myapp notifiers-update --disable <ID>

 $ scalingo -a myapp notifiers-update \
	 --name "My notifier" \
	 --webhook-url https://custom-webhook.com \
	 --send-all-events true \
	 <ID>

 # See also 'notifiers' and 'notifiers-remove'
	`,
		Before: AuthenticateHook,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			var err error

			var active *bool
			if c.IsSet("enable") {
				tmpActive := true
				active = &tmpActive
			} else if c.IsSet("disable") {
				tmpActive := false
				active = &tmpActive
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

			params := scalingo.NotifierParams{
				Active:         active,
				Name:           c.String("name"),
				SendAllEvents:  sendAllEvents,
				SelectedEvents: strings.Split(c.String("events"), " "),

				// Type data options
				PhoneNumber: c.String("phone"),
				Email:       c.String("email"),
				WebhookURL:  c.String("webhook-url"),
			}
			if len(c.Args()) >= 1 {
				err = notifiers.Update(currentApp, c.Args()[0], params)
			} else {
				cli.ShowCommandHelp(c, "notifiers-update")
			}
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "notifiers-update")
		},
	}

	NotifiersRemoveCommand = cli.Command{
		Name:     "notifiers-remove",
		Category: "Notifiers",
		Flags:    []cli.Flag{appFlag},
		Usage:    "Remove an existing notifier from your app",
		Description: `Remove an existing notifier from your app:
    $ scalingo -a myapp notifier-remove <ID>

		# See also 'notifiers' and 'notifiers-add'
`,
		Before: AuthenticateHook,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			var err error
			if len(c.Args()) == 1 {
				err = notifiers.Destroy(currentApp, c.Args()[0])
			} else {
				cli.ShowCommandHelp(c, "notifiers-remove")
			}
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "notifiers-remove")
			autocomplete.NotifiersRemoveAutoComplete(c)
		},
	}
)
