package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
)

var (
	CreateCommand = cli.Command{
		Name:        "create",
		Aliases:     []string{"c"},
		Category:    "Global",
		Description: "Create a new app:\n   Example:\n     'scalingo create mynewapp'\n     'scalingo create mynewapp --remote \"staging\"'",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "remote", Value: "scalingo", Usage: "Remote to add to your current git repository"},
			&cli.StringFlag{Name: "buildpack", Value: "", Usage: "URL to a custom buildpack that Scalingo should use to build your application"},
		},
		Usage: "Create a new app",
		Action: func(c *cli.Context) error {
			if c.Args().Len() != 1 {
				_ = cli.ShowCommandHelp(c, "create")
				return nil
			}
			err := apps.Create(c.Context, c.Args().First(), c.String("remote"), c.String("buildpack"))
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "create")
		},
	}
)
