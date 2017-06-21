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
		Usage:    "List enabled notifiers",
		Flags:    []cli.Flag{appFlag},
		Description: ` List all notifiers enabled for your app:
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
			cli.StringFlag{Name: "active, act", Value: "", Usage: "Use this to enable or disable the notifier. Default: true"},
			cli.StringFlag{Name: "platform, p", Value: "", Usage: "The notifier platform"},
			cli.StringFlag{Name: "name, n", Value: "", Usage: "Name of the notifier"},
			cli.BoolFlag{Name: "send-all-events, sa", Usage: "If true the notifier will send all events. Default: false"},
			cli.StringFlag{Name: "events, e", Value: "", Usage: "List of selected events. Default: []"},
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
   --webhook-url https://hooks.slack.com/services/1234 \
   --events "deployment stop_app"

 $ scalingo -a myapp notifiers-add \
   --platform webhook \
   --name "My notifier" \
   --webhook-url https://custom-webhook.com \
   --send-all-events true

 # See also 'notifiers' and 'notifiers-remove'
`,
		Before: AuthenticateHook,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)

			var active bool
			if c.IsSet("active") {
				active = c.Bool("active")
			} else {
				active = true
			}

			params := scalingo.NotifierCreateParams{
				Active:         active,
				Name:           c.String("name"),
				SendAllEvents:  c.Bool("send-all-events"),
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

	// 	NotificationsUpdateCommand = cli.Command{
	// 		Name:     "notifications-update",
	// 		Category: "Notifications",
	// 		Flags:    []cli.Flag{appFlag},
	// 		Usage:    "Update the url of a notification",
	// 		Description: ` Update the url of a notification:
	//     $ scalingo -a myapp notifications-update <ID> <new-url>
	//
	// 		# See also 'notifications' and 'notifications-add'
	// `,
	// 		Before: AuthenticateHook,
	// 		Action: func(c *cli.Context) {
	// 			currentApp := appdetect.CurrentApp(c)
	// 			var err error
	// 			if len(c.Args()) == 2 {
	// 				err = notifications.Update(currentApp, c.Args()[0], c.Args()[1])
	// 			} else {
	// 				cli.ShowCommandHelp(c, "notifications-update")
	// 			}
	// 			if err != nil {
	// 				errorQuit(err)
	// 			}
	// 		},
	// 		BashComplete: func(c *cli.Context) {
	// 			autocomplete.CmdFlagsAutoComplete(c, "notifications-update")
	// 		},
	// 	}

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
