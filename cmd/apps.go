package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
)

var (
	appsCommand = cli.Command{
		Name:        "apps",
		Category:    "Global",
		Description: "List your apps and give some details about them",
		Flags:       []cli.Flag{&cli.StringFlag{Name: "project", Usage: "Filter apps by project. The filter needs to use the format <ownerUsername>/<projectName>"}},
		Usage:       "List your apps",
		Action: func(c *cli.Context) error {
			if err := apps.List(c.Context, c.String("project")); err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "apps")
		},
	}

	appsInfoCommand = cli.Command{
		Name:     "apps-info",
		Category: "App Management",
		Flags:    []cli.Flag{&appFlag},
		Usage:    "Display the application information",
		Description: CommandDescription{
			Description: "Display various application information such as the force HTTPS status, the stack configured, sticky sessions, etc.",
			Examples:    []string{"scalingo apps-info --app my-app"},
		}.Render(),

		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			if err := apps.Info(c.Context, currentApp); err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "apps-info")
		},
	}
)
