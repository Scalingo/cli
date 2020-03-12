package cmd

import (
	"github.com/urfave/cli"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
)

var (
	CreateCommand = cli.Command{
		Name:        "create",
		ShortName:   "c",
		Category:    "Global",
		Description: "Create a new app:\n   Example:\n     'scalingo create mynewapp'\n     'scalingo create mynewapp --remote \"staging\"'",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "remote", Value: "scalingo", Usage: "Remote to add to your current git repository", EnvVar: ""},
			cli.StringFlag{Name: "buildpack", Value: "", Usage: "URL to a custom buildpack that Scalingo should use to build your application", EnvVar: ""},
		},
		Usage: "Create a new app",
		Action: func(c *cli.Context) {
			if len(c.Args()) != 1 {
				_ = cli.ShowCommandHelp(c, "create")
				return
			}
			err := apps.Create(c.Args()[0], c.String("remote"), c.String("buildpack"))
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "create")
		},
	}
)
