package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/git"
	"github.com/Scalingo/cli/utils"
)

var (
	gitSetup = cli.Command{
		Name:     "git-setup",
		Category: "Git",
		Usage:    "Configure the Git remote for this application",
		Flags: []cli.Flag{&appFlag,
			&cli.StringFlag{
				Name:    "remote",
				Aliases: []string{"r"},
				Value:   "scalingo",
				Usage:   "Specify the remote name"},
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Replace remote url even if remote already exist"},
		},
		Description: CommandDescription{
			Description: "Add a Git remote to the current folder",
			Examples:    []string{"scalingo --app my-app git-setup --remote scalingo-staging"},
		}.Render(),

		Action: func(c *cli.Context) error {
			if c.Args().Len() != 0 {
				cli.ShowCommandHelp(c, "git-setup")
				return nil
			}

			currentApp := detect.CurrentApp(c)
			utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeContainers)

			err := git.Setup(c.Context, currentApp, git.SetupParams{
				RemoteName:     detect.RemoteNameFromFlags(c),
				ForcePutRemote: c.Bool("force"),
			})
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "git-setup")
		},
	}
	gitShow = cli.Command{
		Name:     "git-show",
		Category: "Git",
		Usage:    "Display the Git remote URL for this application",
		Flags:    []cli.Flag{&appFlag},
		Description: CommandDescription{
			Description: "Display the Git remote URL for this application",
			Examples:    []string{"scalingo --app my-app git-show"},
		}.Render(),

		Action: func(c *cli.Context) error {
			if c.Args().Len() != 0 {
				cli.ShowCommandHelp(c, "git-show")
				return nil
			}
			currentApp := detect.CurrentApp(c)
			utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeContainers)

			err := git.Show(c.Context, currentApp)
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "git-show")
		},
	}
)
