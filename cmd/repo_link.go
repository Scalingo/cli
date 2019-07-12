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
		Usage:    "Show repo link linked with your app",
		Flags:    []cli.Flag{appFlag},
		Description: ` Show repo link linked with your application:
	$ scalingo -a myapp repo-link

		# See also 'repo-link-create', 'repo-link-update' and 'repo-link-delete'`,
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

	RepoLinkCreateCommand = cli.Command{
		Name:     "repo-link-create",
		Category: "Repo Link",
		Flags:    []cli.Flag{appFlag},
		Usage:    "Create a repo link between your scm integration and your app",
		Description: ` Create a repo link between your scm integration and your application:
	$ scalingo -a myapp repo-link-create <integration-name> <repo-http-url>
									   OR
	$ scalingo -a myapp repo-link-create <integration-uuid> <repo-http-url>

	Examples:
	$ scalingo -a test-app repo-link-create gitlab https://gitlab.com/gitlab-org/gitlab-ce

		# See also 'repo-link', 'repo-link-update' and 'repo-link-delete'`,
		Action: func(c *cli.Context) {
			var err error

			currentApp := appdetect.CurrentApp(c)

			if len(c.Args()) == 2 {
				err = repo_link.Create(currentApp, c.Args()[0], c.Args()[1])
			} else {
				_ = cli.ShowCommandHelp(c, "repo-link-create")
			}

			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "repo-link-create")
		},
	}

	RepoLinkDeleteCommand = cli.Command{
		Name:     "repo-link-delete",
		Category: "Repo Link",
		Flags:    []cli.Flag{appFlag},
		Usage:    "Delete a repo link linked with your app",
		Description: `Delete a repo link linked with your app:

	$ scalingo -a myapp repo-link-delete repo-link-uuid

		# See also 'repo-link', 'repo-link-create', 'repo-link-update'`,
		Action: func(c *cli.Context) {
			var err error

			currentApp := appdetect.CurrentApp(c)

			if len(c.Args()) == 1 {
				err = repo_link.Delete(currentApp, c.Args()[0])
			} else {
				_ = cli.ShowCommandHelp(c, "repo-link-delete")
			}

			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "repo-link-delete")
		},
	}
)
