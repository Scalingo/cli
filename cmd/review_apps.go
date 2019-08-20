package cmd

import (
	"github.com/urfave/cli"

	"github.com/Scalingo/cli/review_apps"

	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/cmd/autocomplete"
)

var (
	reviewAppsShowCommand = cli.Command{
		Name:     "review-apps",
		Category: "Review Apps",
		Flags:    []cli.Flag{appFlag},
		Usage:    "See review apps of parent application",
		Description: `See review apps of parent application:

	$ scalingo -a myapp review-apps`,
		Action: func(c *cli.Context) {
			var err error

			currentApp := appdetect.CurrentApp(c)

			if len(c.Args()) == 0 {
				err = review_apps.Show(currentApp)
			} else {
				_ = cli.ShowCommandHelp(c, "review-apps")
			}

			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "review-apps")
		},
	}
)
