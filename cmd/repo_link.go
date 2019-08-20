package cmd

import (
	"github.com/urfave/cli"

	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/repo_link"
	"github.com/Scalingo/go-scalingo"
)

var (
	repoLinkShowCommand = cli.Command{
		Name:     "repo-link",
		Category: "Repo Link",
		Usage:    "Show repo link linked with your app",
		Flags:    []cli.Flag{appFlag},
		Description: ` Show repo link linked with your application:
	$ scalingo --app my-app repo-link

		# See also 'repo-link-create', 'repo-link-update', 'repo-link-delete', 'repo-link-manual-deploy', 'repo-link-manual-review-app'`,
		Action: func(c *cli.Context) {
			if len(c.Args()) != 0 {
				cli.ShowCommandHelp(c, "repo-link")
				return
			}

			currentApp := appdetect.CurrentApp(c)
			err := repo_link.Show(currentApp)
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "repo-link")
		},
	}

	repoLinkCreateCommand = cli.Command{
		Name:     "repo-link-create",
		Category: "Repo Link",
		Flags: []cli.Flag{
			appFlag,
			cli.StringFlag{Name: "branch", Usage: "Branch used in auto-deploy"},
			cli.StringFlag{Name: "auto-deploy", Usage: "Enable auto-deploy of application after each branch change"},
			cli.StringFlag{Name: "deploy-review-apps", Usage: "Enable auto-deploy of review app when new pull request is opened"},
			cli.StringFlag{Name: "delete-on-close", Usage: "Enable auto-delete of review apps when pull request is closed"},
			cli.StringFlag{Name: "hours-before-delete-on-close", Usage: "Given time delay of auto-delete of review apps when pull request is closed"},
			cli.StringFlag{Name: "delete-on-stale", Usage: "Enable auto-delete of review apps when no deploy/commits is happen"},
			cli.StringFlag{Name: "hours-before-delete-on-stale", Usage: "Given time delay of auto-delete of review apps when no deploy/commits is happen"},
		},
		Usage: "Create a repo link between your scm integration and your app",
		Description: ` Create a repo link between your scm integration and your application:
	$ scalingo --app my-app repo-link-create <integration-name> <repo-http-url> [options]
									   OR
	$ scalingo --app my-app repo-link-create <integration-uuid> <repo-http-url> [options]

	List of available integrations:
	- github => GitHub.com
	- github-enterprise => GitHub Enterprise (private instance)
	- gitlab => GitLab.com
	- gitlab-self-hosted => GitLab Self-hosted (private instance)

	Examples:
	$ scalingo -a test-app repo-link-create gitlab https://gitlab.com/gitlab-org/gitlab-ce
	$ scalingo -a test-app repo-link-create github-enterprise https://ghe.example.org/test/frontend-app --branch master --auto-deploy true

		# See also 'repo-link', 'repo-link-update', 'repo-link-delete', 'repo-link-manual-deploy' and 'repo-link-manual-review-app'`,
		Action: func(c *cli.Context) {
			if len(c.Args()) != 2 {
				cli.ShowCommandHelp(c, "repo-link-create")
				return
			}

			currentApp := appdetect.CurrentApp(c)
			integrationType := c.Args()[0]
			integrationURL := c.Args()[1]
			branch := c.String("branch")
			autoDeploy := c.Bool("auto-deploy")
			deployReviewApps := c.Bool("deploy-review-apps")
			deleteOnClose := c.Bool("delete-on-close")
			hoursBeforeDeleteOnClose := c.Uint("hours-before-delete-on-close")
			deleteStale := c.Bool("delete-on-stale")
			hoursBeforeDeleteStale := c.Uint("hours-before-delete-on-stale")

			params := scalingo.SCMRepoLinkParams{
				Branch:                   &branch,
				AutoDeployEnabled:        &autoDeploy,
				DeployReviewAppsEnabled:  &deployReviewApps,
				DestroyOnCloseEnabled:    &deleteOnClose,
				HoursBeforeDeleteOnClose: &hoursBeforeDeleteOnClose,
				DestroyStaleEnabled:      &deleteStale,
				HoursBeforeDeleteStale:   &hoursBeforeDeleteStale,
			}

			err := repo_link.Create(currentApp, integrationType, integrationURL, params)
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "repo-link-create")
		},
	}

	repoLinkUpdateCommand = cli.Command{
		Name:     "repo-link-update",
		Category: "Repo Link",
		Flags: []cli.Flag{
			appFlag,
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
	$ scalingo --app my-app repo-link-update [options]

	Examples:
	$ scalingo --app my-app repo-link-update --branch master
	$ scalingo --app my-app repo-link-update --auto-deploy true --branch test --deploy-review-apps true
	$ scalingo --app my-app repo-link-update --delete-on-close true --hours-before-delete-on-close 1
	$ scalingo --app my-app repo-link-update --delete-on-stale true --hours-before-delete-on-stale 2

		# See also 'repo-link', 'repo-link-create', 'repo-link-delete', 'repo-link-manual-deploy' and 'repo-link-manual-review-app'`,
		Action: func(c *cli.Context) {
			if c.NumFlags() == 0 || len(c.Args()) != 0 {
				cli.ShowCommandHelp(c, "repo-link-update")
				return
			}

			currentApp := appdetect.CurrentApp(c)
			params, err := repo_link.CheckAndFillParams(c, currentApp)
			if err != nil {
				errorQuit(err)
			}

			err = repo_link.Update(currentApp, *params)
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "repo-link-update")
		},
	}

	repoLinkDeleteCommand = cli.Command{
		Name:     "repo-link-delete",
		Category: "Repo Link",
		Flags:    []cli.Flag{appFlag},
		Usage:    "Delete repo link linked with your app",
		Description: `Delete repo link linked with your app:

	$ scalingo --app my-app repo-link-delete

		# See also 'repo-link', 'repo-link-create', 'repo-link-update', 'repo-link-manual-deploy' and 'repo-link-manual-review-app'`,
		Action: func(c *cli.Context) {
			if len(c.Args()) != 0 {
				cli.ShowCommandHelp(c, "repo-link-delete")
				return
			}

			currentApp := appdetect.CurrentApp(c)
			err := repo_link.Delete(currentApp)
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "repo-link-delete")
		},
	}

	repoLinkManualDeployCommand = cli.Command{
		Name:     "repo-link-manual-deploy",
		Category: "Repo Link",
		Flags:    []cli.Flag{appFlag},
		Usage:    "Trigger a manual deployment from the state of the branch specified",
		Description: `Trigger a manual deployment from the state of the branch specified:

	$ scalingo --app my-app repo-link-manual-deploy mybranch

		# See also 'repo-link', 'repo-link-create', 'repo-link-update', 'repo-link-delete' and 'repo-link-manual-review-app'`,
		Action: func(c *cli.Context) {
			if len(c.Args()) != 1 {
				cli.ShowCommandHelp(c, "repo-link-manual-deploy")
				return
			}

			currentApp := appdetect.CurrentApp(c)
			branchName := c.Args()[0]
			err := repo_link.ManualDeploy(currentApp, branchName)
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "repo-link-manual-deploy")
		},
	}

	repoLinkManualReviewAppCommand = cli.Command{
		Name:     "repo-link-manual-review-app",
		Category: "Repo Link",
		Flags:    []cli.Flag{appFlag},
		Usage:    "Trigger a manual deployment of review app from the state of the pull/merge request id specified",
		Description: `Trigger a manual deployment of review app from the state of the pull/merge request id specified:

	$ scalingo --app my-app repo-link-manual-review-app pull-request-id (for GitHub and GitHub Enterprise)
	$ scalingo --app my-app repo-link-manual-review-app merge-request-id (for GitLab and GitLab self-hosted)

	Example:
	$ scalingo --app my-app repo-link-manual-review-app 42

		# See also 'repo-link', 'repo-link-create', 'repo-link-update', 'repo-link-delete' and 'repo-link-manual-deploy'`,
		Action: func(c *cli.Context) {
			if len(c.Args()) != 1 {
				cli.ShowCommandHelp(c, "repo-link-manual-review-app")
				return
			}

			currentApp := appdetect.CurrentApp(c)
			pullRequestID := c.Args()[0]
			err := repo_link.ManualReviewApp(currentApp, pullRequestID)
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "repo-link-manual-review-app")
		},
	}
)
