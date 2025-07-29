package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/projects"
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
)
