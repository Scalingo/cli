package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/collaborators"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v8"
)

var (
	CollaboratorsListCommand = cli.Command{
		Name:        "collaborators",
		Category:    "Collaborators",
		Usage:       "List the collaborators of an application",
		Flags:       []cli.Flag{&appFlag},
		Description: "List all the collaborators of an application and display information about them",
		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() != 0 {
				err := cli.ShowCommandHelp(ctx, c, "collaborators")
				if err != nil {
					errorQuit(ctx, err)
				}

				return nil
			}

			err := collaborators.List(ctx, currentApp)
			if err != nil {
				errorQuit(ctx, err)
			}

			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "collaborators")
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
		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() != 1 {
				err := cli.ShowCommandHelp(ctx, c, "collaborators-add")
				if err != nil {
					errorQuit(ctx, err)
				}

				return nil
			}

			utils.CheckForConsent(ctx, currentApp, utils.ConsentTypeContainers)
			err := collaborators.Add(ctx, currentApp, scalingo.CollaboratorAddParams{Email: c.Args().First(), IsLimited: false})
			if err != nil {
				errorQuit(ctx, err)
			}

			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "collaborators-add")
			_ = autocomplete.CollaboratorsAddAutoComplete(ctx, c)
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
		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() != 1 {
				err := cli.ShowCommandHelp(ctx, c, "collaborators-remove")
				if err != nil {
					errorQuit(ctx, err)
				}

				return nil
			}

			utils.CheckForConsent(ctx, currentApp, utils.ConsentTypeContainers)
			err := collaborators.Remove(ctx, currentApp, c.Args().First())
			if err != nil {
				errorQuit(ctx, err)
			}

			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "collaborators-remove")
			_ = autocomplete.CollaboratorsGenericListAutoComplete(ctx, c)
		},
	}

	CollaboratorsUpdateCommand = cli.Command{
		Name:      "collaborators-update",
		Category:  "Collaborators",
		Usage:     "Update a collaborator from an application",
		ArgsUsage: "email",
		Flags: []cli.Flag{
			&appFlag,
			&cli.BoolFlag{Name: "limited", Value: false, Usage: "Set to true if you want to update this collaborator with the limited role"},
		},
		Description: CommandDescription{
			Description: "Update a collaborator from an application, allowing to mark it as limited collaborator or not",
			Examples:    []string{"scalingo --app my-app collaborators-update --limited=true user@example.com"},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() != 1 {
				err := cli.ShowCommandHelp(ctx, c, "collaborators-update")
				if err != nil {
					errorQuit(ctx, err)
				}

				return nil
			}

			utils.CheckForConsent(ctx, currentApp, utils.ConsentTypeContainers)
			err := collaborators.Update(ctx, currentApp, c.Args().First(), scalingo.CollaboratorUpdateParams{IsLimited: c.Bool("limited")})
			if err != nil {
				errorQuit(ctx, err)
			}

			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "collaborators-update")
			_ = autocomplete.CollaboratorsGenericListAutoComplete(ctx, c)
		},
	}
)
