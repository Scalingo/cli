package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"gopkg.in/errgo.v1"

	"github.com/AlecAivazis/survey/v2"
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/integrationlink"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/scmintegrations"
	"github.com/Scalingo/go-scalingo/v6"
	"github.com/Scalingo/go-scalingo/v6/http"
	scalingoerrors "github.com/Scalingo/go-utils/errors"
)

var (
	reviewAppsFromForksSecurityWarning = "Only allow automatic review apps deployments from forks if you trust the owners of those forks, as this could lead to security issues. More info here: https://doc.scalingo.com/platform/app/review-apps#addons-collaborators-and-environment-variables"

	integrationLinkShowCommand = cli.Command{
		Name:     "integration-link",
		Category: "Integration Link",
		Usage:    "Show integration link of your app",
		Flags:    []cli.Flag{&appFlag},
		Description: CommandDescription{
			Description: "Show integration link of your app",
			Examples:    []string{"scalingo --app my-app integration-link"},
			SeeAlso:     []string{"integration-link-create", "integration-link-update", "integration-link-delete", "integration-link-manual-deploy", "integration-link-manual-review-app"},
		}.Render(),
		Action: func(c *cli.Context) error {
			if c.Args().Len() != 0 {
				cli.ShowCommandHelp(c, "integration-link")
				return nil
			}

			currentApp := detect.CurrentApp(c)
			err := integrationlink.Show(c.Context, currentApp)
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "integration-link")
		},
	}

	integrationLinkCreateCommand = cli.Command{
		Name:     "integration-link-create",
		Category: "Integration Link",
		Flags: []cli.Flag{
			&appFlag,
			&cli.StringFlag{Name: "branch", Usage: "Branch used in auto deploy"},
			&cli.BoolFlag{Name: "auto-deploy", Usage: "Enable auto deploy of application after each branch change"},
			&cli.BoolFlag{Name: "deploy-review-apps", Usage: "Enable auto deploy of review apps when new pull request is opened"},
			&cli.BoolFlag{Name: "destroy-on-close", Usage: "Auto destroy review apps when pull request is closed"},
			&cli.BoolFlag{Name: "no-auto-deploy", Usage: "Disable auto deploy of application after each branch change"},
			&cli.BoolFlag{Name: "no-deploy-review-apps", Usage: "Disable auto deploy of review app when new pull request is opened"},
			&cli.BoolFlag{Name: "no-destroy-on-close", Usage: "Auto destroy review apps when pull request is closed"},
			&cli.BoolFlag{Name: "allow-review-apps-from-forks", Usage: "Enable auto deploy of review apps when new pull request is opened from a fork"},
			&cli.BoolFlag{Name: "aware-of-security-risks", Usage: "Bypass the security warning about allowing automatic review app creation from forks"},
			&cli.BoolFlag{Name: "no-allow-review-apps-from-forks", Usage: "Disable auto deploy of review apps when new pull request is opened from a fork"},
			&cli.UintFlag{Name: "hours-before-destroy-on-close", Usage: "Time delay before auto destroying a review app when pull request is closed"},
			&cli.BoolFlag{Name: "destroy-on-stale", Usage: "Auto destroy review apps when no deploy/commits has happened"},
			&cli.BoolFlag{Name: "no-destroy-on-stale", Usage: "Auto destroy review apps when no deploy/commits has happened"},
			&cli.UintFlag{Name: "hours-before-destroy-on-stale", Usage: "Time delay before auto destroying a review app when no deploy/commits has happened"},
		},
		Usage:     "Link your Scalingo application to an integration",
		ArgsUsage: "repository-url",
		Description: CommandDescription{
			Description: `Link your Scalingo application to an integration

List of available integrations:
- github => GitHub.com
- github-enterprise => GitHub Enterprise (private instance)
- gitlab => GitLab.com
- gitlab-self-hosted => GitLab Self-hosted (private instance)			
`,
			Examples: []string{
				"scalingo --app my-app integration-link-create https://gitlab.com/gitlab-org/gitlab-ce",
				"scalingo --app my-app integration-link-create --branch master --auto-deploy --deploy-review-apps --no-allow-review-apps-from-forks https://ghe.example.org/test/frontend-app",
			},
			SeeAlso: []string{"integration-link", "integration-link-update", "integration-link-delete", "integration-link-manual-deploy", "integration-link-manual-review-app"},
		}.Render(),
		Action: func(c *cli.Context) error {
			if c.Args().Len() != 1 {
				cli.ShowCommandHelp(c, "integration-link-create")
				return nil
			}

			currentApp := detect.CurrentApp(c)
			integrationURL := c.Args().First()
			integrationURLParsed, err := url.Parse(integrationURL)
			if err != nil {
				errorQuit(errgo.Notef(err, "error parsing the repository url"))
			}
			// If the customer forgot to specify the scheme, we automatically prefix with https://
			if integrationURLParsed.Scheme == "" {
				integrationURL = fmt.Sprintf("https://%s", integrationURL)
			}

			integrationType, err := scmintegrations.GetTypeFromURL(c.Context, integrationURL)
			if err != nil {
				if scalingoerrors.ErrgoRoot(err) == scmintegrations.ErrNotFound {
					// If no integration matches the given URL, display a helpful status
					// message
					switch integrationURLParsed.Host {
					case "github.com":
						io.Error("No GitHub integration found, please follow this URL to add it:")
						io.Errorf("%s/users/github/link\n", config.C.ScalingoAuthUrl)
					case "gitlab.com":
						io.Error("No GitLab integration found, please follow this URL to add it:")
						io.Errorf("%s/users/gitlab/link\n", config.C.ScalingoAuthUrl)
					default:
						io.Errorf("No integration found for URL %s.\n", integrationURL)
						io.Errorf("Please run the command:\n")
						io.Errorf("scalingo integrations-add gitlab-self-hosted|github-enterprise --url %s://%s --token <personal-access-token>\n", integrationURLParsed.Scheme, integrationURLParsed.Host)
					}
					os.Exit(1)
				}
				errorQuit(err)
			}

			var params scalingo.SCMRepoLinkCreateParams
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

				allowReviewAppsFromForks := c.Bool("allow-review-apps-from-forks")
				noAllowReviewAppsFromForks := c.Bool("no-allow-review-apps-from-forks")

				if allowReviewAppsFromForks && noAllowReviewAppsFromForks {
					errorQuit(errors.New("cannot define both allow-review-apps-from-forks and no-allow-review-apps-from-forks"))
				}

				awareOfSecurityRisks := c.Bool("aware-of-security-risks")

				if deployReviewApps && allowReviewAppsFromForks && !awareOfSecurityRisks {
					allowReviewAppsFromForks, err = askForConfirmation(reviewAppsFromForksSecurityWarning)
					if err != nil {
						errorQuit(err)
					}
				}

				params = scalingo.SCMRepoLinkCreateParams{
					Branch:                            &branch,
					AutoDeployEnabled:                 &autoDeploy,
					DeployReviewAppsEnabled:           &deployReviewApps,
					DestroyOnCloseEnabled:             &destroyOnClose,
					HoursBeforeDeleteOnClose:          &hoursBeforeDestroyOnClose,
					DestroyStaleEnabled:               &destroyOnStale,
					HoursBeforeDeleteStale:            &hoursBeforeDestroyOnStale,
					AutomaticCreationFromForksAllowed: &allowReviewAppsFromForks,
				}
			}

			err = integrationlink.Create(c.Context, currentApp, integrationType, integrationURL, params)
			if err != nil {
				scerr, ok := scalingoerrors.ErrgoRoot(err).(*http.RequestFailedError)
				if ok {
					if scerr.Code == 404 {
						io.Error("Fail to create SCM repository integration: the repository has not been found")
						io.Errorf("Check %v exists and you have the correct permissions\n", integrationURL)
						if integrationType == scalingo.SCMGithubType || integrationType == scalingo.SCMGithubEnterpriseType {
							io.Error("https://doc.scalingo.com/platform/deployment/deploy-with-github")
						} else if integrationType == scalingo.SCMGitlabType || integrationType == scalingo.SCMGitlabSelfHostedType {
							io.Error("https://doc.scalingo.com/platform/deployment/deploy-with-gitlab")
						}
					} else if scerr.Code == 403 {
						io.Error("Fail to create SCM repository integration: the SCM API returned a 401 error.")
						io.Error("Did you revoked the Scalingo token from your profile? In such situation, you may want to remove and re-create the SCM integration.")
						io.Error("")
						io.Errorf("The complete error message from the SCM API is: %s\n", scerr.APIError)
					} else {
						errorQuit(err)
					}
				} else {
					errorQuit(err)
				}
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "integration-link-create")
		},
	}

	integrationLinkUpdateCommand = cli.Command{
		Name:     "integration-link-update",
		Category: "Integration Link",
		Flags: []cli.Flag{
			&appFlag,
			&cli.StringFlag{Name: "branch", Usage: "Branch used in auto deploy"},
			&cli.BoolFlag{Name: "auto-deploy", Usage: "Enable auto deploy of application after each branch change"},
			&cli.BoolFlag{Name: "no-auto-deploy", Usage: "Disable auto deploy of application after each branch change"},
			&cli.BoolFlag{Name: "deploy-review-apps", Usage: "Enable auto deploy of review app when new pull request is opened"},
			&cli.BoolFlag{Name: "no-deploy-review-apps", Usage: "Disable auto deploy of review app when new pull request is opened"},
			&cli.BoolFlag{Name: "allow-review-apps-from-forks", Usage: "Enable auto deploy of review apps when new pull request is opened from a fork"},
			&cli.BoolFlag{Name: "aware-of-security-risks", Usage: "Bypass the security warning about allowing automatic review app creation from forks"},
			&cli.BoolFlag{Name: "no-allow-review-apps-from-forks", Usage: "Disable auto deploy of review apps when new pull request is opened from a fork"},
			&cli.BoolFlag{Name: "destroy-on-close", Usage: "Auto destroy review apps when pull request is closed"},
			&cli.BoolFlag{Name: "no-destroy-on-close", Usage: "Auto destroy review apps when pull request is closed"},
			&cli.UintFlag{Name: "hours-before-destroy-on-close", Usage: "Time delay before auto destroying a review app when pull request is closed"},
			&cli.BoolFlag{Name: "destroy-on-stale", Usage: "Auto destroy review apps when no deploy/commits has happened"},
			&cli.BoolFlag{Name: "no-destroy-on-stale", Usage: "Auto destroy review apps when no deploy/commits has happened"},
			&cli.StringFlag{Name: "hours-before-destroy-on-stale", Usage: "Time delay before auto destroying a review app when no deploy/commits has happened"},
		},
		Usage: "Update the integration link parameters",
		Description: CommandDescription{
			Description: "Update the integration link parameters",
			Examples: []string{
				"scalingo --app my-app integration-link-update --branch master",
				"scalingo --app my-app integration-link-update --auto-deploy --branch test --deploy-review-apps",
				"scalingo --app my-app integration-link-update --destroy-on-close --hours-before-destroy-on-close 1",
				"scalingo --app my-app integration-link-update --destroy-on-stale --hours-before-destroy-on-stale 2",
			},
			SeeAlso: []string{"integration-link", "integration-link-create", "integration-link-delete", "integration-link-manual-deploy", "integration-link-manual-review-app"},
		}.Render(),
		Action: func(c *cli.Context) error {
			if c.NumFlags() == 0 || c.Args().Len() != 0 {
				cli.ShowCommandHelp(c, "integration-link-update")
				return nil
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

			allowReviewAppsFromForks := c.Bool("allow-review-apps-from-forks")
			noAllowReviewAppsFromForks := c.Bool("no-allow-review-apps-from-forks")

			if allowReviewAppsFromForks && noAllowReviewAppsFromForks {
				errorQuit(errors.New("cannot define both allow-review-apps-from-forks and no-allow-review-apps-from-forks"))
			}

			awareOfSecurityRisks := c.Bool("aware-of-security-risks")

			if allowReviewAppsFromForks && !awareOfSecurityRisks {
				stillAllowed, err := askForConfirmation(reviewAppsFromForksSecurityWarning)
				if err != nil {
					errorQuit(err)
				}
				err = c.Set("allow-review-apps-from-forks", strconv.FormatBool(stillAllowed))
				c.Value("allow-review-apps-from-forks")
				if err != nil {
					errorQuit(errgo.Notef(err, "error updating if review apps creation from forks are allowed"))
				}
			}

			currentApp := detect.CurrentApp(c)
			params := integrationlink.CheckAndFillParams(c)

			err := integrationlink.Update(c.Context, currentApp, *params)
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "integration-link-update")
		},
	}

	integrationLinkDeleteCommand = cli.Command{
		Name:     "integration-link-delete",
		Category: "Integration Link",
		Flags:    []cli.Flag{&appFlag},
		Usage:    "Delete the link between your Scalingo application and the integration",
		Description: CommandDescription{
			Description: "Delete the link between your Scalingo application and the integration",
			Examples:    []string{"scalingo --app my-app integration-link-delete"},
			SeeAlso:     []string{"integration-link", "integration-link-create", "integration-link-update", "integration-link-manual-deploy", "integration-link-manual-review-app"},
		}.Render(),
		Action: func(c *cli.Context) error {
			if c.Args().Len() != 0 {
				cli.ShowCommandHelp(c, "integration-link-delete")
				return nil
			}

			currentApp := detect.CurrentApp(c)
			err := integrationlink.Delete(c.Context, currentApp)
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "integration-link-delete")
		},
	}

	integrationLinkManualDeployCommand = cli.Command{
		Name:      "integration-link-manual-deploy",
		Category:  "Integration Link",
		Flags:     []cli.Flag{&appFlag},
		Usage:     "Trigger a manual deployment of the specified branch",
		ArgsUsage: "branch",
		Description: CommandDescription{
			Description: "Trigger a manual deployment of the specified branch",
			Examples:    []string{"scalingo --app my-app integration-link-manual-deploy mybranch"},
			SeeAlso:     []string{"integration-link", "integration-link-create", "integration-link-update", "integration-link-delete", "integration-link-manual-review-app"},
		}.Render(),
		Action: func(c *cli.Context) error {
			if c.Args().Len() != 1 {
				cli.ShowCommandHelp(c, "integration-link-manual-deploy")
				return nil
			}

			currentApp := detect.CurrentApp(c)
			branchName := c.Args().First()
			err := integrationlink.ManualDeploy(c.Context, currentApp, branchName)
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "integration-link-manual-deploy")
		},
	}

	integrationLinkManualReviewAppCommand = cli.Command{
		Name:      "integration-link-manual-review-app",
		Category:  "Integration Link",
		Flags:     []cli.Flag{&appFlag},
		Usage:     "Trigger a review app creation of the pull/merge request ID specified",
		ArgsUsage: "request-id",
		Description: CommandDescription{
			Description: `Trigger a review app creation of the pull/merge request ID specified:

   $ scalingo --app my-app integration-link-manual-review-app pull-request-id (for GitHub and GitHub Enterprise)
   $ scalingo --app my-app integration-link-manual-review-app merge-request-id (for GitLab and GitLab self-hosted)
`,
			Examples: []string{"scalingo --app my-app integration-link-manual-review-app 42"},
			SeeAlso:  []string{"integration-link", "integration-link-create", "integration-link-update", "integration-link-delete", "integration-link-manual-deploy"},
		}.Render(),
		Action: func(c *cli.Context) error {
			if c.Args().Len() != 1 {
				cli.ShowCommandHelp(c, "integration-link-manual-review-app")
				return nil
			}

			currentApp := detect.CurrentApp(c)
			pullRequestID := c.Args().First()

			err := integrationlink.ManualReviewApp(c.Context, currentApp, pullRequestID)
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "integration-link-manual-review-app")
		},
	}
)

func interactiveCreate() (scalingo.SCMRepoLinkCreateParams, error) {
	var params scalingo.SCMRepoLinkCreateParams
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
		return params, errgo.Notef(err, "error enquiring about branch and automatic review apps deployment")
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
		Default: destroyOnClose,
	}, &destroyOnClose, nil)
	if err != nil {
		return params, errgo.Notef(err, "error enquiring about destroy on close")
	}
	params.DestroyOnCloseEnabled = &destroyOnClose
	if destroyOnClose {
		answerHoursBeforeDestroyOnClose := "0"
		err = survey.AskOne(&survey.Input{
			Message: "Hours before automatically destroying the review apps:",
			Default: "0",
		}, &answerHoursBeforeDestroyOnClose, survey.WithValidator(validateHoursBeforeDelete))
		if err != nil {
			return params, errgo.Notef(err, "error enquiring about review apps destroy delay")
		}
		hoursBeforeDestroyOnClose64, _ := strconv.ParseUint(answerHoursBeforeDestroyOnClose, 10, 32)
		hoursBeforeDestroyOnClose := uint(hoursBeforeDestroyOnClose64)
		params.HoursBeforeDeleteOnClose = &hoursBeforeDestroyOnClose
	}

	destroyOnStale := false
	err = survey.AskOne(&survey.Confirm{
		Message: "Automatically destroy review apps after some time without deploy/commits:",
		Default: destroyOnStale,
	}, &destroyOnStale, nil)
	if err != nil {
		return params, errgo.Notef(err, "error enquiring about stale review apps destroy")
	}
	params.DestroyStaleEnabled = &destroyOnStale
	if destroyOnStale {
		answerHoursBeforeDestroyOnStale := "0"
		err = survey.AskOne(&survey.Input{
			Message: "Hours before automatically destroying the review apps:",
			Default: "0",
		}, &answerHoursBeforeDestroyOnStale, survey.WithValidator(validateHoursBeforeDelete))
		if err != nil {
			return params, errgo.Notef(err, "error enquiring about stale review apps destroy")
		}
		hoursBeforeDestroyOnStale64, _ := strconv.ParseUint(answerHoursBeforeDestroyOnStale, 10, 32)
		hoursBeforeDestroyOnStale := uint(hoursBeforeDestroyOnStale64)
		params.HoursBeforeDeleteStale = &hoursBeforeDestroyOnStale
	}

	io.Warning(reviewAppsFromForksSecurityWarning)
	var forksAllowed bool

	err = survey.AskOne(&survey.Confirm{
		Message: "Allow review apps to be created from forks:",
		Default: false,
	}, &forksAllowed, nil)

	if err != nil {
		return params, errgo.Notef(err, "error enquiring about automatic review apps creation from forks")
	}
	params.AutomaticCreationFromForksAllowed = &forksAllowed

	return params, nil
}

func validateHoursBeforeDelete(ans interface{}) error {
	str, ok := ans.(string)
	if !ok {
		return errors.New("must be a string")
	}
	i, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return errgo.Notef(err, "error parsing hours")
	}
	if i < 0 {
		return errors.New("must be positive")
	}
	return nil
}

func askForConfirmation(message string) (bool, error) {
	io.Warning(message)
	var confirmed bool

	err := survey.AskOne(&survey.Confirm{
		Message: "Are your sure?",
		Default: false,
	}, &confirmed, nil)

	if err != nil {
		return false, err
	}

	return confirmed, nil
}
