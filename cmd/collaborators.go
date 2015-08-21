package cmd

import (
	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/collaborators"
)

var (
	CollaboratorsListCommand = cli.Command{
		Name:        "collaborators",
		Category:    "Collaborators",
		Usage:       "List the collaborators of an application",
		Description: "List all the collaborator of an application and display information about them.",
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			if len(c.Args()) != 0 {
				cli.ShowCommandHelp(c, "collaborators")
			} else {
				err := collaborators.List(currentApp)
				if err != nil {
					errorQuit(err)
				}
			}
		},
	}

	CollaboratorsAddCommand = cli.Command{
		Name:        "collaborators-add",
		Category:    "Collaborators",
		Usage:       "Invite someone to work on an application",
		Description: "Invite someone to collaborate on an application, an invitation will be sent to the given email\n scalingo -a myapp collaborators-add user@example.com",
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			if len(c.Args()) != 1 {
				cli.ShowCommandHelp(c, "collaborators-add")
			} else {
				err := collaborators.Add(currentApp, c.Args()[0])
				if err != nil {
					errorQuit(err)
				}
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CollaboratorsAddAutoComplete(c)
		},
	}

	CollaboratorsRemoveCommand = cli.Command{
		Name:        "collaborators-remove",
		Category:    "Collaborators",
		Usage:       "Revoke permission to collaborate on an application",
		Description: "Revoke the invitation of collaboration to the given email\n    scalingo -a myapp collaborators-remove user@example.com",
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
			if len(c.Args()) != 1 {
				cli.ShowCommandHelp(c, "collaborators-remove")
			} else {
				err := collaborators.Remove(currentApp, c.Args()[0])
				if err != nil {
					errorQuit(err)
				}
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CollaboratorsRemoveAutoComplete(c)
		},
	}
)
