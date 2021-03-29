package cmd

import (
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/urfave/cli"
)

var (
	oneOffStopCommand = cli.Command{
		Name:     "one-off-stop",
		Category: "App Management",
		Usage:    "Stop a running one-off container",
		Flags:    []cli.Flag{appFlag},
		Description: `Stop a running one-off container
	Example
	  'scalingo --app my-app one-off-stop one-off-1234'`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			if len(c.Args()) != 1 {
				cli.ShowCommandHelp(c, "one-off-stop")
				return
			}

			err := apps.OneOffStop(currentApp, c.Args()[0])
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "one-off-stop")
		},
	}
)
