package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/projects"
	"github.com/Scalingo/go-scalingo/v8"
)

var (
	projectsListCommand = cli.Command{
		Name:        "projects",
		Category:    "Projects",
		Usage:       "List the projects that you own",
		Description: "List all the projects of which you are an owner",
		Action: func(c *cli.Context) error {

			err := projects.List(c.Context)
			if err != nil {
				errorQuit(c.Context, err)
			}

			return nil
		},
		BashComplete: func(c *cli.Context) {
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
				"scalingo projects-add my-awesome-project			# Create a new project",
				"scalingo projects-add --default my-awesome-project # Create a new project that is the default one",
			},
		}.Render(),
		Action: func(c *cli.Context) error {
			err := projects.Add(c.Context, scalingo.ProjectAddParams{Name: c.Args().First(), Default: c.Bool("default")})
			if err != nil {
				errorQuit(c.Context, err)
			}

			return nil
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "projects-add")
		},
	}
)
