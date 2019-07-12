package cmd

import (
	"github.com/urfave/cli"

	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/repo_link"
	"github.com/Scalingo/go-scalingo"
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

	RepoLinkUpdateCommand = cli.Command{
		Name:     "repo-link-update",
		Category: "Repo Link",
		Flags: []cli.Flag{appFlag,
			cli.StringFlag{Name: "branch", Usage: "Branch", EnvVar: ""},
			cli.StringFlag{Name: "auto-deploy", Usage: "Enable auto-deploy", EnvVar: ""},
			cli.StringFlag{Name: "deploy-review-apps", Usage: "Enable auto-deploy of review apps", EnvVar: ""},
			cli.StringFlag{Name: "delete-on-close", Usage: "Enable delete review apps when pull request is closed", EnvVar: ""},
			cli.StringFlag{Name: "hours-before-delete-on-close", Usage: "", EnvVar: ""},
			cli.StringFlag{Name: "delete-stale", Usage: "", EnvVar: ""},
			cli.StringFlag{Name: "hours-before-delete-stale", Usage: "", EnvVar: ""},
		},
		Usage: "Update the repo link linked with your app",
		Description: ` Update the repo link linked with your application:
	$ scalingo -a myapp repo-link-update --<options>

	Examples:
	$ scalingo -a myapp repo-link-update --branch master
	$ scalingo -a myapp repo-link-update --auto-deploy true --branch test --deploy-review-apps true
	$ scalingo -a myapp repo-link-update --delete-on-close true --hours-before-delete-on-close 1
	$ scalingo -a myapp repo-link-update --delete-stale true --hours-before-delete-stale 2

		# See also 'repo-link', 'repo-link-create' and 'repo-link-delete'`,
		Action: func(c *cli.Context) {
			var err error

			currentApp := appdetect.CurrentApp(c)

			params := scalingo.ScmRepoLinkParams{
				Branch:                   c.String("branch"),
				AutoDeployEnabled:        c.Bool("auto-deploy"),
				DeployReviewAppsEnabled:  c.Bool("deploy-review-apps"),
				DestroyOnCloseEnabled:    c.Bool("delete-on-close"),
				HoursBeforeDeleteOnClose: c.Uint("hours-before-delete-on-close"),
				DestroyStaleEnabled:      c.Bool("delete-stale"),
				HoursBeforeDeleteStale:   c.Uint("hours-before-delete-stale"),
			}

			if len(c.Args()) == 0 {
				err = repo_link.Update(currentApp, params)
			} else {
				_ = cli.ShowCommandHelp(c, "repo-link-update")
			}

			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "repo-link-update")
		},
	}

	RepoLinkDeleteCommand = cli.Command{
		Name:     "repo-link-delete",
		Category: "Repo Link",
		Flags:    []cli.Flag{appFlag},
		Usage:    "Delete repo link linked with your app",
		Description: `Delete repo link linked with your app:

	$ scalingo -a myapp repo-link-delete

		# See also 'repo-link', 'repo-link-create', 'repo-link-update'`,
		Action: func(c *cli.Context) {
			var err error

			currentApp := appdetect.CurrentApp(c)

			if len(c.Args()) == 0 {
				err = repo_link.Delete(currentApp)
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
