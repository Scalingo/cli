package cmd

import (
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/notifiers"
	"github.com/Scalingo/codegangsta-cli"
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
				cli.ShowCommandHelp(c, "notifications")
			}

			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "notifications")
		},
	}

// 	NotificationsAddCommand = cli.Command{
// 		Name:     "notifications-add",
// 		Category: "Notifications",
// 		Flags:    []cli.Flag{appFlag},
// 		Usage:    "Enable a notification for your application",
// 		Description: ` Enable a notification for your application:
//     $ scalingo -a myapp notifications-add <webhook-url>
//
// 		# See also 'notifications' and 'notifications-remove'
// `,
// 		Before: AuthenticateHook,
// 		Action: func(c *cli.Context) {
// 			currentApp := appdetect.CurrentApp(c)
// 			var err error
// 			if len(c.Args()) == 1 {
// 				err = notifications.Provision(currentApp, c.Args()[0])
// 			} else {
// 				cli.ShowCommandHelp(c, "notifications-add")
// 			}
// 			if err != nil {
// 				errorQuit(err)
// 			}
// 		},
// 		BashComplete: func(c *cli.Context) {
// 			autocomplete.CmdFlagsAutoComplete(c, "notifications-add")
// 		},
// 	}
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
