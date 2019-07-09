package cmd

import (
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/git"
	"github.com/urfave/cli"
)

var (
	gitSetup = cli.Command{
		Name:     "git-setup",
		Category: "Git",
		Usage:    "Configure the Git remote for this application",
		Flags: []cli.Flag{appFlag,
			cli.StringFlag{
				Name: "remote, r", Value: "scalingo",
				Usage: "Specify the remote name"},
		},
		Description: `Add a Git remote to the current folder.

		Example
		  scalingo --app my-app git-setup --remote scalingo-staging
		`,
		Action: func(c *cli.Context) {
			if len(c.Args()) != 0 {
				cli.ShowCommandHelp(c, "git-setup")
				return
			}

			err := git.Setup(appdetect.CurrentApp(c), git.SetupParams{
				RemoteName: c.String("r"),
			})
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "git-setup")
		},
	}
	gitShow = cli.Command{
		Name:     "git-show",
		Category: "Git",
		Usage:    "Display the Git remote URL for this application",
		Flags:    []cli.Flag{appFlag},
		Description: `Display the Git remote URL for this application.

		Example
		  scalingo --app my-app git-show
		`,
		Action: func(c *cli.Context) {
			if len(c.Args()) != 0 {
				cli.ShowCommandHelp(c, "git-show")
				return
			}

			err := git.Show(appdetect.CurrentApp(c))
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "git-show")
		},
	}
)
