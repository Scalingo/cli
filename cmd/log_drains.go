package cmd

import (
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/log_drains"
	"github.com/urfave/cli"
)

var (
	logDrainsListCommand = cli.Command{
		Name:     "log-drains",
		Category: "Log drains",
		Flags:    []cli.Flag{appFlag},
		Usage:    "List the log drains of an application",
		Description: `List all the log drains of an application:

    $ scalingo --app my-app log-drains`,

		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			if len(c.Args()) != 0 {
				cli.ShowCommandHelp(c, "log-drains")
				return
			}

			err := log_drains.List(currentApp)
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "log-drains")
		},
	}
)
