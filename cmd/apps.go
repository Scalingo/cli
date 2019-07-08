package cmd

import (
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/urfave/cli"
)

var (
	appsCommand = cli.Command{
		Name:        "apps",
		Category:    "Global",
		Description: "List your apps and give some details about them",
		Usage:       "List your apps",
		Action: func(c *cli.Context) {
			if err := apps.List(); err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "apps")
		},
	}

	appsInfoCommand = cli.Command{
		Name:     "apps-info",
		Category: "App Management",
		Flags:    []cli.Flag{appFlag},
		Usage:    "Display the application information",
		Description: `Display various application information such as the force HTTPS status, the stack configured, sticky sessions, etc.

		Example:
			scalingo apps-info --app my-app
`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			if err := apps.Info(currentApp); err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "apps-info")
		},
	}
)
