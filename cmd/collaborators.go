package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/collaborators"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/utils"
)

var (
	CollaboratorsListCommand = cli.Command{
		Name:        "collaborators",
		Category:    "Collaborators",
		Usage:       "List the collaborators of an application",
		Flags:       []cli.Flag{&appFlag},
		Description: "List all the collaborators of an application and display information about them",
		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() != 0 {
				cli.ShowCommandHelp(c, "collaborators")
			} else {
				err := collaborators.List(c.Context, currentApp)
				if err != nil {
					errorQuit(err)
				}
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "collaborators")
		},
	}

	CollaboratorsAddCommand = cli.Command{
		Name:      "collaborators-add",
		Category:  "Collaborators",
		Usage:     "Invite someone to work on an application",
		ArgsUsage: "email",
		Flags:     []cli.Flag{&appFlag},
		Description: CommandDescription{
			Description: "Invite someone to collaborate on an application, an invitation will be sent to the given email",
			Examples:    []string{"scalingo --app my-app collaborators-add user@example.com"},
		}.Render(),
		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() != 1 {
				cli.ShowCommandHelp(c, "collaborators-add")
			} else {
				utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeContainers)
				err := collaborators.Add(c.Context, currentApp, c.Args().First())
				if err != nil {
					errorQuit(err)
				}
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "collaborators-add")
			autocomplete.CollaboratorsAddAutoComplete(c)
		},
	}

	CollaboratorsRemoveCommand = cli.Command{
		Name:      "collaborators-remove",
		Category:  "Collaborators",
		Usage:     "Revoke permission to collaborate on an application",
		ArgsUsage: "email",
		Flags:     []cli.Flag{&appFlag},
		Description: CommandDescription{
			Description: "Revoke the invitation of collaboration to the given email",
			Examples:    []string{"scalingo -a myapp collaborators-remove user@example.com"},
		}.Render(),
		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() != 1 {
				cli.ShowCommandHelp(c, "collaborators-remove")
			} else {
				utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeContainers)
				err := collaborators.Remove(c.Context, currentApp, c.Args().First())
				if err != nil {
					errorQuit(err)
				}
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "collaborators-remove")
			autocomplete.CollaboratorsRemoveAutoComplete(c)
		},
	}
)
