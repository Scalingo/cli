package cmd

import (
	"github.com/urfave/cli"

	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/repo_link"
)

var (
	RepoLinkShowCommand = cli.Command{
		Name:     "repo-link",
		Category: "Repo Link",
		Usage:    "Show repo link linked with your app",
		Flags:    []cli.Flag{appFlag},
		Description: ` Show repo link linked with your application:
	$ scalingo -a myapp repo-link

		# See also 'repo-link-create', 'repo-link-update', 'repo-link-delete', 'repo-link-manual-deploy', 'repo-link-manual-review-app'`,
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

		# See also 'repo-link', 'repo-link-update', 'repo-link-delete', 'repo-link-manual-deploy' and 'repo-link-manual-review-app'`,
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
			cli.StringFlag{Name: "branch", Usage: "Branch used in auto-deploy"},
			cli.StringFlag{Name: "auto-deploy", Usage: "Enable auto-deploy of application after each branch change"},
			cli.StringFlag{Name: "deploy-review-apps", Usage: "Enable auto-deploy of review app when new pull request is opened"},
			cli.StringFlag{Name: "delete-on-close", Usage: "Enable auto-delete of review apps when pull request is closed"},
			cli.StringFlag{Name: "hours-before-delete-on-close", Usage: "Given time delay of auto-delete of review apps when pull request is closed"},
			cli.StringFlag{Name: "delete-on-stale", Usage: "Enable auto-delete of review apps when no deploy/commits is happen"},
			cli.StringFlag{Name: "hours-before-delete-on-stale", Usage: "Given time delay of auto-delete of review apps when no deploy/commits is happen"},
		},
		Usage: "Update the repo link linked with your app",
		Description: ` Update the repo link linked with your application:
	$ scalingo -a myapp repo-link-update --<options>

	Examples:
	$ scalingo -a myapp repo-link-update --branch master
	$ scalingo -a myapp repo-link-update --auto-deploy true --branch test --deploy-review-apps true
	$ scalingo -a myapp repo-link-update --delete-on-close true --hours-before-delete-on-close 1
	$ scalingo -a myapp repo-link-update --delete-on-stale true --hours-before-delete-on-stale 2

		# See also 'repo-link', 'repo-link-create', 'repo-link-delete', 'repo-link-manual-deploy' and 'repo-link-manual-review-app'`,
		Action: func(c *cli.Context) {
			var err error

			currentApp := appdetect.CurrentApp(c)

			if c.NumFlags() == 0 {
				_ = cli.ShowCommandHelp(c, "repo-link-update")
				return
			}

			params, err := repo_link.CheckAndFillParams(c, currentApp)
			if err != nil {
				errorQuit(err)
			}

			if len(c.Args()) == 0 {
				err = repo_link.Update(currentApp, *params)
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

		# See also 'repo-link', 'repo-link-create', 'repo-link-update', 'repo-link-manual-deploy' and 'repo-link-manual-review-app'`,
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

	RepoLinkManualDeployCommand = cli.Command{
		Name:     "repo-link-manual-deploy",
		Category: "Repo Link",
		Flags:    []cli.Flag{appFlag},
		Usage:    "Trigger a manual deployment from the state of the branch specified.",
		Description: `Trigger a manual deployment from the state of the branch specified:

	$ scalingo -a myapp repo-link-manual-deploy mybranch

		# See also 'repo-link', 'repo-link-create', 'repo-link-update', 'repo-link-delete' and 'repo-link-manual-review-app'`,
		Action: func(c *cli.Context) {
			var err error

			currentApp := appdetect.CurrentApp(c)

			if len(c.Args()) == 1 {
				err = repo_link.ManualDeploy(currentApp, c.Args()[0])
			} else {
				_ = cli.ShowCommandHelp(c, "repo-link-manual-deploy")
			}

			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "repo-link-manual-deploy")
		},
	}

	RepoLinkManualReviewAppCommand = cli.Command{
		Name:     "repo-link-manual-review-app",
		Category: "Repo Link",
		Flags:    []cli.Flag{appFlag},
		Usage:    "Trigger a manual deployment of review app from the state of the pull/merge request id specified.",
		Description: `Trigger a manual deployment of review app from the state of the pull/merge request id specified:

	$ scalingo -a myapp repo-link-manual-review-app pull-request-id (for GitHub and GitHub Enterprise)
	$ scalingo -a myapp repo-link-manual-review-app merge-request-id (for GitLab and GitLab self-hosted)

	Example:
	$ scalingo -a myapp repo-link-manual-review-app 42

		# See also 'repo-link', 'repo-link-create', 'repo-link-update', 'repo-link-delete' and 'repo-link-manual-deploy'`,
		Action: func(c *cli.Context) {
			var err error

			currentApp := appdetect.CurrentApp(c)

			if len(c.Args()) == 1 {
				err = repo_link.ManualReviewApp(currentApp, c.Args()[0])
			} else {
				_ = cli.ShowCommandHelp(c, "repo-link-manual-review-app")
			}

			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "repo-link-manual-review-app")
		},
	}

	RepoLinkReviewAppsCommand = cli.Command{
		Name:     "repo-link-review-apps",
		Category: "Repo Link",
		Flags:    []cli.Flag{appFlag},
		Usage:    "See reviews apps for th",
		Description: `Trigger a manual deployment of review app from the state of the pull/merge request id specified:

	$ scalingo -a myapp repo-link-manual-review-app pull-request-id (for GitHub and GitHub Enterprise)
	$ scalingo -a myapp repo-link-manual-review-app merge-request-id (for GitLab and GitLab self-hosted)

	Example:
	$ scalingo -a myapp repo-link-manual-review-app 42

		# See also 'repo-link', 'repo-link-create', 'repo-link-update', 'repo-link-delete' and 'repo-link-manual-deploy'`,
		Action: func(c *cli.Context) {
			var err error

			currentApp := appdetect.CurrentApp(c)

			if len(c.Args()) == 1 {
				err = repo_link.ManualReviewApp(currentApp, c.Args()[0])
			} else {
				_ = cli.ShowCommandHelp(c, "repo-link-manual-review-app")
			}

			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "repo-link-manual-review-app")
		},
	}
)
