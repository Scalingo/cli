package cmd

import (
	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/cli/notifications"
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/cmd/autocomplete"
)


var (
	NotificationsListCommand = cli.Command{
		Name:     "notifications",
		Category: "Notifications",
		Usage:    "List enabled notifications",
		Flags:    []cli.Flag{appFlag},
		Description: ` List all notifications enabled for your app:
    $ scalingo -a myapp notifications

		# See also 'notifications-add' and 'notifications-remove'
`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			var err error
			if len(c.Args()) == 0 {
				err = notifications.List(currentApp)
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
	NotificationsAddCommand = cli.Command{
		Name:     "notifications-add",
		Category: "Notifications",
		Flags:    []cli.Flag{appFlag},
		Usage:    "Enable a notification for your application",
		Description: ` Enable a notification for your application:
    $ scalingo -a myapp notifications-add <webhook-url>

		# See also 'notifications-list' and 'notifications-remove'
`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			var err error
			if len(c.Args()) == 1 {
				err = notifications.Provision(currentApp, c.Args()[0])
			} else {
				cli.ShowCommandHelp(c, "notifications-add")
			}
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "notifications-add")
		},
	}
	NotificationsRemoveCommand = cli.Command{
		Name:     "notifications-remove",
		Category: "Notifications",
		Flags:    []cli.Flag{appFlag},
		Usage:    "Remove an existing notification from your app",
		Description: ` Remove an existing notification from your app:
    $ scalingo -a myapp notifications-remove <ID>

		# See also 'notifications' and 'notifications-add'
`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			var err error
			if len(c.Args()) == 1 {
				err = notifications.Destroy(currentApp, c.Args()[0])
			} else {
				cli.ShowCommandHelp(c, "notifications-remove")
			}
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "notifications-remove")
		},
	}
)
