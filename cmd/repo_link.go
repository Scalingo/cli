package cmd

import (
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/repo_link"
	"github.com/urfave/cli"
)

var (
	RepoLinkShowCommand = cli.Command{
		Name:     "repo-link",
		Category: "Repo Link",
		Usage:    "Show repo link linked to your app",
		Flags:    []cli.Flag{appFlag},
		Description: ` Show repo link linked to your app:
    $ scalingo -a myapp repo-link

		# See also 'repo-link-add', 'repo-link-update' and 'repo-link-remove'`,
		Action: func(c *cli.Context) {
			var err error

			currentApp := appdetect.CurrentApp(c)

			if len(c.Args()) == 0 {
				err = repo_link.Show(currentApp)
			} else {
				_ = cli.ShowCommandHelp(c, "repo-link")
			}

			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "repo-link")
		},
	}
)
