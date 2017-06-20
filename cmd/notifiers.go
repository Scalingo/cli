package cmd

import (
	"fmt"

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
			cli.StringFlag{Name: "platform, p", Value: "", Usage: "The notifier platform"},
			cli.StringFlag{Name: "name, n", Value: "", Usage: "Name of the notifier"},
			cli.BoolFlag{Name: "send-all-events, sa", Usage: "Should the notifier send all events or not"},
			cli.StringFlag{Name: "events, e", Value: "", Usage: "List of selected events"},
			cli.StringFlag{Name: "phone", Value: "", Usage: "A phone number"},
			cli.StringFlag{Name: "email", Value: "", Usage: "An email"},
			cli.StringFlag{Name: "webhook-url, u", Value: "", Usage: "A webhook url"},
		},
		Usage: "Add a notifier for your application",
		Description: ` Add a notifier for your application:
    $ scalingo -a myapp notifiers-add <webhook-url>

		# See also 'notifiers' and 'notifiers-remove'
`,
		Before: AuthenticateHook,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			params := scalingo.NotifierCreateParams{
				Name:           c.String("name"),
				SendAllEvents:  c.Bool("send-all-events"),
				SelectedEvents: []string{c.String("events")}, //strings.Split(c.String("events"), " "),

				// Type data options
				PhoneNumber: c.String("phone"),
				Email:       c.String("email"),
				WebhookURL:  c.String("webhook-url"),
			}

			var err error
			fmt.Println(len(c.Args()))
			// if len(c.Args()) != 1 {
			err = notifiers.Provision(currentApp, c.String("platform"), params)
			// } else {
			// 	cli.ShowCommandHelp(c, "notifiers-add")
			// }
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "notifications-add")
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
// 	NotificationsRemoveCommand = cli.Command{
// 		Name:     "notifications-remove",
// 		Category: "Notifications",
// 		Flags:    []cli.Flag{appFlag},
// 		Usage:    "Remove an existing notification from your app",
// 		Description: ` Remove an existing notification from your app:
//     $ scalingo -a myapp notifications-remove <ID>
//
// 		# See also 'notifications' and 'notifications-add'
// `,
// 		Before: AuthenticateHook,
// 		Action: func(c *cli.Context) {
// 			currentApp := appdetect.CurrentApp(c)
// 			var err error
// 			if len(c.Args()) == 1 {
// 				err = notifications.Destroy(currentApp, c.Args()[0])
// 			} else {
// 				cli.ShowCommandHelp(c, "notifications-remove")
// 			}
// 			if err != nil {
// 				errorQuit(err)
// 			}
// 		},
// 		BashComplete: func(c *cli.Context) {
// 			autocomplete.CmdFlagsAutoComplete(c, "notifications-remove")
// 			autocomplete.NotificationsRemoveAutoComplete(c)
// 		},
// 	}
)
