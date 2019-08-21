package cmd

import (
	"errors"

	"github.com/urfave/cli"

	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/repolink"
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
			err := repolink.Show(currentApp)
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
			cli.StringFlag{Name: "branch", Usage: "Branch used in auto deploy"},
			cli.BoolFlag{Name: "auto-deploy", Usage: "Enable auto deploy of application after each branch change"},
			cli.BoolFlag{Name: "deploy-review-apps", Usage: "Enable auto deploy of review app when new pull request is opened"},
			cli.BoolFlag{Name: "destroy-on-close", Usage: "Auto destroy review apps when pull request is closed"},
			cli.BoolFlag{Name: "no-auto-deploy", Usage: "Enable auto deploy of application after each branch change"},
			cli.BoolFlag{Name: "no-deploy-review-apps", Usage: "Enable auto deploy of review app when new pull request is opened"},
			cli.BoolFlag{Name: "no-destroy-on-close", Usage: "Auto destroy review apps when pull request is closed"},
			cli.UintFlag{Name: "hours-before-destroy-on-close", Usage: "Time delay before auto destroying a review app when pull request is closed"},
			cli.BoolFlag{Name: "destroy-on-stale", Usage: "Auto destroy review apps when no deploy/commits has happened"},
			cli.BoolFlag{Name: "no-destroy-on-stale", Usage: "Auto destroy review apps when no deploy/commits has happened"},
			cli.UintFlag{Name: "hours-before-destroy-on-stale", Usage: "Time delay before auto destroying a review app when no deploy/commits has happened"},
		},
		Usage: "Create a repo link between your scm integration and your app",
		Description: ` Create a repo link between your scm integration and your application:
	$ scalingo --app my-app repo-link-create <integration-name> <repo-url> [options]
									   OR
	$ scalingo --app my-app repo-link-create <integration-uuid> <repo-url> [options]

	List of available integrations:
	- github => GitHub.com
	- github-enterprise => GitHub Enterprise (private instance)
	- gitlab => GitLab.com
	- gitlab-self-hosted => GitLab Self-hosted (private instance)

	Examples:
	$ scalingo -a test-app repo-link-create gitlab https://gitlab.com/gitlab-org/gitlab-ce
	$ scalingo -a test-app repo-link-create github-enterprise https://ghe.example.org/test/frontend-app --branch master --auto-deploy

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
			noAutoDeploy := c.Bool("no-auto-deploy")
			if autoDeploy && noAutoDeploy {
				errorQuit(errors.New("cannot define both auto-deploy and no-auto-deploy"))
			}

			deployReviewApps := c.Bool("deploy-review-apps")
			noDeployReviewApps := c.Bool("no-deploy-review-apps")
			if deployReviewApps && noDeployReviewApps {
				errorQuit(errors.New("cannot define both deploy-review-apps and no-deploy-review-apps"))
			}

			destroyOnClose := c.Bool("destroy-on-close")
			noDestroyOnClose := c.Bool("no-destroy-on-close")
			if destroyOnClose && noDestroyOnClose {
				errorQuit(errors.New("cannot define both destroy-on-close and no-destroy-on-close"))
			}
			hoursBeforeDestroyOnClose := c.Uint("hours-before-destroy-on-close")

			destroyOnStale := c.Bool("destroy-on-stale")
			noDestroyOnStale := c.Bool("no-destroy-on-stale")
			if destroyOnStale && noDestroyOnStale {
				errorQuit(errors.New("cannot define both destroy-on-stale and no-destroy-on-stale"))
			}
			hoursBeforeDestroyOnStale := c.Uint("hours-before-destroy-on-stale")

			params := scalingo.SCMRepoLinkParams{
				Branch:                   &branch,
				AutoDeployEnabled:        &autoDeploy,
				DeployReviewAppsEnabled:  &deployReviewApps,
				DestroyOnCloseEnabled:    &destroyOnClose,
				HoursBeforeDeleteOnClose: &hoursBeforeDestroyOnClose,
				DestroyStaleEnabled:      &destroyOnStale,
				HoursBeforeDeleteStale:   &hoursBeforeDestroyOnStale,
			}

			err := repolink.Create(currentApp, integrationType, integrationURL, params)
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
			cli.StringFlag{Name: "branch", Usage: "Branch used in auto deploy"},
			cli.BoolFlag{Name: "auto-deploy", Usage: "Enable auto deploy of application after each branch change"},
			cli.BoolFlag{Name: "no-auto-deploy", Usage: "Enable auto deploy of application after each branch change"},
			cli.BoolFlag{Name: "deploy-review-apps", Usage: "Enable auto deploy of review app when new pull request is opened"},
			cli.BoolFlag{Name: "no-deploy-review-apps", Usage: "Enable auto deploy of review app when new pull request is opened"},
			cli.BoolFlag{Name: "destroy-on-close", Usage: "Auto destroy review apps when pull request is closed"},
			cli.BoolFlag{Name: "no-destroy-on-close", Usage: "Auto destroy review apps when pull request is closed"},
			cli.UintFlag{Name: "hours-before-destroy-on-close", Usage: "Time delay before auto destroying a review app when pull request is closed"},
			cli.BoolFlag{Name: "destroy-on-stale", Usage: "Auto destroy review apps when no deploy/commits has happened"},
			cli.BoolFlag{Name: "no-destroy-on-stale", Usage: "Auto destroy review apps when no deploy/commits has happened"},
			cli.StringFlag{Name: "hours-before-destroy-on-stale", Usage: "Time delay before auto destroying a review app when no deploy/commits has happened"},
		},
		Usage: "Update the repo link linked with your app",
		Description: ` Update the repo link linked with your application:
	$ scalingo --app my-app repo-link-update [options]

	Examples:
	$ scalingo --app my-app repo-link-update --branch master
	$ scalingo --app my-app repo-link-update --auto-deploy --branch test --deploy-review-apps
	$ scalingo --app my-app repo-link-update --destroy-on-close --hours-before-destroy-on-close 1
	$ scalingo --app my-app repo-link-update --destroy-on-stale --hours-before-destroy-on-stale 2

		# See also 'repo-link', 'repo-link-create', 'repo-link-delete', 'repo-link-manual-deploy' and 'repo-link-manual-review-app'`,
		Action: func(c *cli.Context) {
			if c.NumFlags() == 0 || len(c.Args()) != 0 {
				cli.ShowCommandHelp(c, "repo-link-update")
				return
			}
			autoDeploy := c.Bool("auto-deploy")
			noAutoDeploy := c.Bool("no-auto-deploy")
			if autoDeploy && noAutoDeploy {
				errorQuit(errors.New("cannot define both auto-deploy and no-auto-deploy"))
			}
			deployReviewApps := c.Bool("deploy-review-apps")
			noDeployReviewApps := c.Bool("no-deploy-review-apps")
			if deployReviewApps && noDeployReviewApps {
				errorQuit(errors.New("cannot define both deploy-review-apps and no-deploy-review-apps"))
			}
			destroyOnClose := c.Bool("destroy-on-close")
			noDestroyOnClose := c.Bool("no-destroy-on-close")
			if destroyOnClose && noDestroyOnClose {
				errorQuit(errors.New("cannot define both destroy-on-close and no-destroy-on-close"))
			}
			destroyOnStale := c.Bool("destroy-on-stale")
			noDestroyOnStale := c.Bool("no-destroy-on-stale")
			if destroyOnStale && noDestroyOnStale {
				errorQuit(errors.New("cannot define both destroy-on-stale and no-destroy-on-stale"))
			}

			currentApp := appdetect.CurrentApp(c)
			params, err := repolink.CheckAndFillParams(c, currentApp)
			if err != nil {
				errorQuit(err)
			}

			err = repolink.Update(currentApp, *params)
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
			err := repolink.Delete(currentApp)
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
		Usage:    "Trigger a manual deployment of the specified branch",
		Description: `Trigger a manual deployment of the specified branch:

	$ scalingo --app my-app repo-link-manual-deploy mybranch

		# See also 'repo-link', 'repo-link-create', 'repo-link-update', 'repo-link-delete' and 'repo-link-manual-review-app'`,
		Action: func(c *cli.Context) {
			if len(c.Args()) != 1 {
				cli.ShowCommandHelp(c, "repo-link-manual-deploy")
				return
			}

			currentApp := appdetect.CurrentApp(c)
			branchName := c.Args()[0]
			err := repolink.ManualDeploy(currentApp, branchName)
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
		Usage:    "Trigger a review app creation of the pull/merge request ID specified",
		Description: `Trigger a review app creation of the pull/merge request ID specified:

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
			err := repolink.ManualReviewApp(currentApp, pullRequestID)
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "repo-link-manual-review-app")
		},
	}
)
