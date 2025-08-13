package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

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

		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() != 0 {
				_ = cli.ShowCommandHelp(ctx, c, "git-setup")
				return nil
			}

			currentApp := detect.CurrentApp(c)
			utils.CheckForConsent(ctx, currentApp, utils.ConsentTypeContainers)

			err := git.Setup(ctx, currentApp, git.SetupParams{
				RemoteName:     detect.RemoteNameFromFlags(c),
				ForcePutRemote: c.Bool("force"),
			})
			if err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "git-setup")
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

		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() != 0 {
				_ = cli.ShowCommandHelp(ctx, c, "git-show")
				return nil
			}
			currentApp := detect.CurrentApp(c)
			utils.CheckForConsent(ctx, currentApp, utils.ConsentTypeContainers)

			err := git.Show(ctx, currentApp)
			if err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "git-show")
		},
	}
)
