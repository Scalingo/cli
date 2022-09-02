package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/review_apps"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
)

var (
	reviewAppsShowCommand = cli.Command{
		Name:     "review-apps",
		Category: "Review Apps",
		Flags:    []cli.Flag{&appFlag},
		Usage:    "Show review apps of the parent application",
		Description: `Show review apps of the parent application:

	$ scalingo --app my-app review-apps`,
		Action: func(c *cli.Context) error {
			if c.Args().Len() != 0 {
				cli.ShowCommandHelp(c, "review-apps")
				return nil
			}

			currentApp := detect.CurrentApp(c)
			err := review_apps.Show(currentApp)
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "review-apps")
		},
	}
)
