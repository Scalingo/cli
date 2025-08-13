package cmd

import (
	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/notificationplatforms"
)

var (
	NotificationPlatformListCommand = cli.Command{
		Name:        "notification-platforms",
		Category:    "Notifiers - Global",
		Description: "List all notification platforms you can use with a notifier.",
		Usage:       "List all notification platforms",

		Action: func(ctx context.Context, c *cli.Command) error {
			err := notificationplatforms.List(c.Context)
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "platforms-list")
		},
	}
)
