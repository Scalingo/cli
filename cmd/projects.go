package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/projects"
	"github.com/Scalingo/go-scalingo/v8"
	"github.com/Scalingo/go-utils/errors/v2"
)

var (
	projectsListCommand = cli.Command{
		Name:        "projects",
		Category:    "Projects",
		Usage:       "List the projects that you own",
		Description: "List all the projects of which you are an owner",
		Action: func(ctx context.Context, _ *cli.Command) error {

			err := projects.List(ctx)
			if err != nil {
				errorQuit(ctx, err)
			}

			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "projects")
		},
	}

	projectsAddCommand = cli.Command{
		Name:      "projects-add",
		Category:  "Projects",
		Usage:     "Create a project",
		ArgsUsage: "project-name",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "default", Value: false, Usage: "Set to true if this should be the new default project"},
		},
		Description: CommandDescription{
			Description: "Create a new project, with the capability to make it the default one",
			Examples: []string{
				"scalingo projects-add my-awesome-project           # Create a new project",
				"scalingo projects-add --default my-awesome-project # Create a new project that is the default one",
			},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			err := projects.Add(ctx, scalingo.ProjectAddParams{Name: c.Args().First(), Default: c.Bool("default")})
			if err != nil {
				errorQuit(ctx, err)
			}

			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "projects-add")
		},
	}

	projectsUpdateCommand = cli.Command{
		Name:      "projects-update",
		Category:  "Projects",
		Usage:     "Update a project",
		ArgsUsage: "project-id",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "name", Value: "", Usage: "Change the project name"},
			&cli.BoolFlag{Name: "default", Value: false, Usage: "Set to `true` if you want to make this project the default one. " +
				"It cannot change from `true` to `false`. To change the default project, update an existing project to be the new default one, or create a new default project."},
		},
		Description: CommandDescription{
			Description: "Update a project, allowing to change the name and setting it as the default one",
			Examples: []string{
				"scalingo projects-update --default prj-00000000-0000-0000-0000-000000000000               # Make the project the default one",
				"scalingo projects-update --name new-project-name prj-00000000-0000-0000-0000-000000000000 # Change the project name",
			},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			projectID := c.Args().First()
			if projectID == "" {
				errorQuitWithHelpMessage(ctx, errors.New(ctx, "missing project ID parameter"), c, "projects-update")
			}

			var newProjectNamePtr *string
			// Only set a new project name if a new name was provided
			if newProjectName := c.String("name"); newProjectName != "" {
				newProjectNamePtr = &newProjectName
			}

			var defaultPtr *bool
			// Only allow "default" to be true, otherwise do nothing
			def := c.Bool("default")
			if def {
				defaultPtr = &def
			}

			if newProjectNamePtr == nil && !def {
				errorQuitWithHelpMessage(ctx, errors.New(ctx, "no parameters were submitted"), c, "projects-update")
			}

			err := projects.Update(ctx, projectID, scalingo.ProjectUpdateParams{Name: newProjectNamePtr, Default: defaultPtr})
			if err != nil {
				errorQuit(ctx, err)
			}

			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "projects-update")
			_ = autocomplete.ProjectsGenericListAutoComplete(ctx)
		},
	}

	projectsRemoveCommand = cli.Command{
		Name:      "projects-remove",
		Category:  "Projects",
		Usage:     "Remove a project",
		ArgsUsage: "project-id",
		Description: CommandDescription{
			Description: "Remove a project, given it is not the default one",
			Examples:    []string{"scalingo projects-remove prj-00000000-0000-0000-0000-000000000000"},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			projectID := c.Args().First()
			if projectID == "" {
				errorQuitWithHelpMessage(ctx, errors.New(ctx, "missing project ID parameter"), c, "projects-remove")
			}

			err := projects.Remove(ctx, projectID)
			if err != nil {
				errorQuit(ctx, err)
			}

			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "projects-remove")
			_ = autocomplete.ProjectsGenericListAutoComplete(ctx)
		},
	}
)
