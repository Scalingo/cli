package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
)

var (
	CreateCommand = cli.Command{
		Name:      "create",
		Aliases:   []string{"c"},
		Category:  "Global",
		Usage:     "Create a new app",
		ArgsUsage: "app-id",
		Description: CommandDescription{
			Description: "Create a new app",
			Examples: []string{
				"scalingo create mynewapp",
				"scalingo --remote \"staging\" create mynewapp",
			},
		}.Render(),
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "remote", Value: "scalingo", Usage: "Remote to add to your current git repository"},
			&cli.StringFlag{Name: "buildpack", Value: "", Usage: "URL to a custom buildpack that Scalingo should use to build your application"},
		},

		Action: func(c *cli.Context) error {
			if c.Args().Len() != 1 {
				_ = cli.ShowCommandHelp(c, "create")
				return nil
			}

			err := apps.Create(c.Context, c.Args().First(), detect.RemoteNameFromFlags(c), c.String("buildpack"))
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
