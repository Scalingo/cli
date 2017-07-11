package cmd

import (
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/notification_platforms"
	"github.com/urfave/cli"
)

var (
	NotificationPlatformListCommand = cli.Command{
		Name:        "platforms-list",
		Category:    "Notifiers - Global",
		Description: "List all notification platforms you can use with a notifier.",
		Usage:       "List all notification platforms",
		Action: func(c *cli.Context) {
			err := notification_platforms.List()
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "platforms-list")
		},
	}
)
