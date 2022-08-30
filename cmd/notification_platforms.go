package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/notification_platforms"
)

var (
	NotificationPlatformListCommand = cli.Command{
		Name:        "notification-platforms",
		Category:    "Notifiers - Global",
		Description: "List all notification platforms you can use with a notifier.",
		Usage:       "List all notification platforms",
		Action: func(c *cli.Context) error {
			err := notification_platforms.List()
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "platforms-list")
		},
	}
)
