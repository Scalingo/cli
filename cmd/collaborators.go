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
	collaboratorsListCommand = cli.Command{
		Name:        "collaborators",
		Category:    "Collaborators",
		Usage:       "List the collaborators of an application",
		Flags:       []cli.Flag{&appFlag, databaseFlag()},
		Description: "List all the collaborators of an application and display information about them",
		Action: func(ctx context.Context, c *cli.Command) error {
			currentResource := detect.GetCurrentResource(ctx, c)
			if c.Args().Len() != 0 {
				err := cli.ShowCommandHelp(ctx, c, "collaborators")
				if err != nil {
					errorQuit(ctx, err)
				}

				return nil
			}

			err := collaborators.List(ctx, currentResource)
			if err != nil {
				errorQuit(ctx, err)
			}

			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "collaborators")
		},
	}

	collaboratorsAddCommand = cli.Command{
		Name:      "collaborators-add",
		Category:  "Collaborators",
		Usage:     "Invite someone to work on an application",
		ArgsUsage: "email",
		Flags:     []cli.Flag{&appFlag, databaseFlag()},
		Description: CommandDescription{
			Description: "Invite someone to collaborate on an application, an invitation will be sent to the given email",
			Examples:    []string{"scalingo --app my-app collaborators-add user@example.com"},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			currentResource := detect.GetCurrentResource(ctx, c)
			if c.Args().Len() != 1 {
				err := cli.ShowCommandHelp(ctx, c, "collaborators-add")
				if err != nil {
					errorQuit(ctx, err)
				}

				return nil
			}

			utils.CheckForConsent(ctx, currentResource, utils.ConsentTypeContainers)
			err := collaborators.Add(ctx, currentResource, scalingo.CollaboratorAddParams{Email: c.Args().First(), IsLimited: false})
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

	collaboratorsRemoveCommand = cli.Command{
		Name:      "collaborators-remove",
		Category:  "Collaborators",
		Usage:     "Revoke permission to collaborate on an application",
		ArgsUsage: "email",
		Flags:     []cli.Flag{&appFlag, databaseFlag()},
		Description: CommandDescription{
			Description: "Revoke the invitation of collaboration to the given email",
			Examples:    []string{"scalingo -a myapp collaborators-remove user@example.com"},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			currentResource := detect.GetCurrentResource(ctx, c)
			if c.Args().Len() != 1 {
				err := cli.ShowCommandHelp(ctx, c, "collaborators-remove")
				if err != nil {
					errorQuit(ctx, err)
				}

				return nil
			}

			utils.CheckForConsent(ctx, currentResource, utils.ConsentTypeContainers)
			err := collaborators.Remove(ctx, currentResource, c.Args().First())
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

	collaboratorsUpdateCommand = cli.Command{
		Name:      "collaborators-update",
		Category:  "Collaborators",
		Usage:     "Update a collaborator from an application",
		ArgsUsage: "email",
		Flags: []cli.Flag{
			&appFlag,
			databaseFlag(),
			&cli.BoolFlag{Name: "limited", Value: false, Usage: "Set to true if you want to update this collaborator with the limited role"},
		},
		Description: CommandDescription{
			Description: "Update a collaborator from an application, allowing to mark it as limited collaborator or not",
			Examples:    []string{"scalingo --app my-app collaborators-update --limited=true user@example.com"},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			currentResource := detect.GetCurrentResource(ctx, c)
			if c.Args().Len() != 1 {
				err := cli.ShowCommandHelp(ctx, c, "collaborators-update")
				if err != nil {
					errorQuit(ctx, err)
				}

				return nil
			}

			utils.CheckForConsent(ctx, currentResource, utils.ConsentTypeContainers)
			err := collaborators.Update(ctx, currentResource, c.Args().First(), scalingo.CollaboratorUpdateParams{IsLimited: c.Bool("limited")})
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
