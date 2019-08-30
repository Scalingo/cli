package cmd

import (
	"errors"
	"strconv"

	"github.com/AlecAivazis/survey"
	"github.com/urfave/cli"

	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/integrationlink"
	"github.com/Scalingo/cli/scmintegrations"
	"github.com/Scalingo/go-scalingo"
)

var (
	integrationLinkShowCommand = cli.Command{
		Name:     "integration-link",
		Category: "Integration Link",
		Usage:    "Show integration link of your app",
		Flags:    []cli.Flag{appFlag},
		Description: ` Show integration link of your app:
	$ scalingo --app my-app integration-link

		# See also 'integration-link-create', 'integration-link-update', 'integration-link-delete', 'integration-link-manual-deploy', 'integration-link-manual-review-app'`,
		Action: func(c *cli.Context) {
			if len(c.Args()) != 0 {
				cli.ShowCommandHelp(c, "integration-link")
				return
			}

			currentApp := appdetect.CurrentApp(c)
			err := integrationlink.Show(currentApp)
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "integration-link")
		},
	}

	integrationLinkCreateCommand = cli.Command{
		Name:     "integration-link-create",
		Category: "Integration Link",
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
		Usage: "Link your Scalingo application to an integration",
		Description: ` Link your Scalingo application to an integration:
	$ scalingo --app my-app integration-link-create <repository URL> [options]
									   OR
	$ scalingo --app my-app integration-link-create <repository URL> [options]

	List of available integrations:
	- github => GitHub.com
	- github-enterprise => GitHub Enterprise (private instance)
	- gitlab => GitLab.com
	- gitlab-self-hosted => GitLab Self-hosted (private instance)

	Examples:
	$ scalingo --app my-app integration-link-create https://gitlab.com/gitlab-org/gitlab-ce
	$ scalingo --app my-app integration-link-create https://ghe.example.org/test/frontend-app --branch master --auto-deploy

		# See also 'integration-link', 'integration-link-update', 'integration-link-delete', 'integration-link-manual-deploy' and 'integration-link-manual-review-app'`,
		Action: func(c *cli.Context) {
			if len(c.Args()) != 1 {
				cli.ShowCommandHelp(c, "integration-link-create")
				return
			}

			currentApp := appdetect.CurrentApp(c)
			integrationURL := c.Args()[0]
			integrationType, err := scmintegrations.GetTypeFromURL(integrationURL)
			if err != nil {
				errorQuit(err)
			}

			var params scalingo.SCMRepoLinkParams
			if c.NumFlags() == 0 {
				params, err = interactiveCreate()
				if err != nil {
					errorQuit(err)
				}
			} else {
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

				params = scalingo.SCMRepoLinkParams{
					Branch:                   &branch,
					AutoDeployEnabled:        &autoDeploy,
					DeployReviewAppsEnabled:  &deployReviewApps,
					DestroyOnCloseEnabled:    &destroyOnClose,
					HoursBeforeDeleteOnClose: &hoursBeforeDestroyOnClose,
					DestroyStaleEnabled:      &destroyOnStale,
					HoursBeforeDeleteStale:   &hoursBeforeDestroyOnStale,
				}
			}

			err = integrationlink.Create(currentApp, integrationType, integrationURL, params)
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "integration-link-create")
		},
	}

	integrationLinkUpdateCommand = cli.Command{
		Name:     "integration-link-update",
		Category: "Integration Link",
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
		Usage: "Update the integration link parameters",
		Description: ` Update the integration link parameters:
	$ scalingo --app my-app integration-link-update [options]

	Examples:
	$ scalingo --app my-app integration-link-update --branch master
	$ scalingo --app my-app integration-link-update --auto-deploy --branch test --deploy-review-apps
	$ scalingo --app my-app integration-link-update --destroy-on-close --hours-before-destroy-on-close 1
	$ scalingo --app my-app integration-link-update --destroy-on-stale --hours-before-destroy-on-stale 2

		# See also 'integration-link', 'integration-link-create', 'integration-link-delete', 'integration-link-manual-deploy' and 'integration-link-manual-review-app'`,
		Action: func(c *cli.Context) {
			if c.NumFlags() == 0 || len(c.Args()) != 0 {
				cli.ShowCommandHelp(c, "integration-link-update")
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
			params, err := integrationlink.CheckAndFillParams(c, currentApp)
			if err != nil {
				errorQuit(err)
			}

			err = integrationlink.Update(currentApp, *params)
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "integration-link-update")
		},
	}

	integrationLinkDeleteCommand = cli.Command{
		Name:     "integration-link-delete",
		Category: "Integration Link",
		Flags:    []cli.Flag{appFlag},
		Usage:    "Delete the link between your Scalingo application and the integration",
		Description: `Delete the link between your Scalingo application and the integration:

	$ scalingo --app my-app integration-link-delete

		# See also 'integration-link', 'integration-link-create', 'integration-link-update', 'integration-link-manual-deploy' and 'integration-link-manual-review-app'`,
		Action: func(c *cli.Context) {
			if len(c.Args()) != 0 {
				cli.ShowCommandHelp(c, "integration-link-delete")
				return
			}

			currentApp := appdetect.CurrentApp(c)
			err := integrationlink.Delete(currentApp)
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "integration-link-delete")
		},
	}

	integrationLinkManualDeployCommand = cli.Command{
		Name:     "integration-link-manual-deploy",
		Category: "Integration Link",
		Flags:    []cli.Flag{appFlag},
		Usage:    "Trigger a manual deployment of the specified branch",
		Description: `Trigger a manual deployment of the specified branch:

	$ scalingo --app my-app integration-link-manual-deploy mybranch

		# See also 'integration-link', 'integration-link-create', 'integration-link-update', 'integration-link-delete' and 'integration-link-manual-review-app'`,
		Action: func(c *cli.Context) {
			if len(c.Args()) != 1 {
				cli.ShowCommandHelp(c, "integration-link-manual-deploy")
				return
			}

			currentApp := appdetect.CurrentApp(c)
			branchName := c.Args()[0]
			err := integrationlink.ManualDeploy(currentApp, branchName)
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "integration-link-manual-deploy")
		},
	}

	integrationLinkManualReviewAppCommand = cli.Command{
		Name:     "integration-link-manual-review-app",
		Category: "Integration Link",
		Flags:    []cli.Flag{appFlag},
		Usage:    "Trigger a review app creation of the pull/merge request ID specified",
		Description: `Trigger a review app creation of the pull/merge request ID specified:

	$ scalingo --app my-app integration-link-manual-review-app pull-request-id (for GitHub and GitHub Enterprise)
	$ scalingo --app my-app integration-link-manual-review-app merge-request-id (for GitLab and GitLab self-hosted)

	Example:
	$ scalingo --app my-app integration-link-manual-review-app 42

		# See also 'integration-link', 'integration-link-create', 'integration-link-update', 'integration-link-delete' and 'integration-link-manual-deploy'`,
		Action: func(c *cli.Context) {
			if len(c.Args()) != 1 {
				cli.ShowCommandHelp(c, "integration-link-manual-review-app")
				return
			}

			currentApp := appdetect.CurrentApp(c)
			pullRequestID := c.Args()[0]
			err := integrationlink.ManualReviewApp(currentApp, pullRequestID)
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "integration-link-manual-review-app")
		},
	}
)

func interactiveCreate() (scalingo.SCMRepoLinkParams, error) {
	var params scalingo.SCMRepoLinkParams
	if config.C.DisableInteractive {
		return params, errors.New("need at least one integration link parameter")
	}
	qs := []*survey.Question{
		{
			Name:   "branch",
			Prompt: &survey.Input{Message: "Branch to auto-deploy (empty to disable):"},
		},
		{
			Name: "auto-review-apps",
			Prompt: &survey.Confirm{
				Message: "Automatically deploy review apps:",
				Default: false,
			},
		},
	}

	answers := struct {
		Branch         string
		AutoReviewApps bool `survey:"auto-review-apps"`
	}{}
	err := survey.Ask(qs, &answers)
	if err != nil {
		return params, err
	}

	t := true
	if answers.Branch != "" {
		params.Branch = &answers.Branch
		params.AutoDeployEnabled = &t
	}
	if !answers.AutoReviewApps {
		return params, nil
	}

	params.DeployReviewAppsEnabled = &t

	destroyOnClose := true
	err = survey.AskOne(&survey.Confirm{
		Message: "Automatically destroy review apps when the pull/merge request is closed:",
		Default: true,
	}, &destroyOnClose, nil)
	if err != nil {
		return params, err
	}
	params.DestroyOnCloseEnabled = &destroyOnClose
	if destroyOnClose {
		answerHoursBeforeDestroyOnClose := "0"
		err = survey.AskOne(&survey.Input{
			Message: "Hours before automatically destroying the review apps:",
			Default: "0",
		}, &answerHoursBeforeDestroyOnClose, validateHoursBeforeDelete)
		if err != nil {
			return params, err
		}
		hoursBeforeDestroyOnClose64, _ := strconv.ParseUint(answerHoursBeforeDestroyOnClose, 10, 32)
		hoursBeforeDestroyOnClose := uint(hoursBeforeDestroyOnClose64)
		params.HoursBeforeDeleteOnClose = &hoursBeforeDestroyOnClose
	}

	destroyOnStale := true
	err = survey.AskOne(&survey.Confirm{
		Message: "Automatically destroy review apps after some time without deploy/commits:",
		Default: true,
	}, &destroyOnStale, nil)
	if err != nil {
		return params, err
	}
	params.DestroyStaleEnabled = &destroyOnStale
	if destroyOnStale {
		answerHoursBeforeDestroyOnStale := "0"
		err = survey.AskOne(&survey.Input{
			Message: "Hours before automatically destroying the review apps:",
			Default: "0",
		}, &answerHoursBeforeDestroyOnStale, validateHoursBeforeDelete)
		if err != nil {
			return params, err
		}
		hoursBeforeDestroyOnStale64, _ := strconv.ParseUint(answerHoursBeforeDestroyOnStale, 10, 32)
		hoursBeforeDestroyOnStale := uint(hoursBeforeDestroyOnStale64)
		params.HoursBeforeDeleteStale = &hoursBeforeDestroyOnStale
	}
	return params, nil
}

func validateHoursBeforeDelete(ans interface{}) error {
	str, ok := ans.(string)
	if !ok {
		return errors.New("must be a string")
	}
	i, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return err
	}
	if i < 0 {
		return errors.New("must be positive")
	}
	return nil
}
