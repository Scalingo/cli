package cmd

import (
	"net/url"

	"github.com/urfave/cli"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/scm_integrations"
	"github.com/Scalingo/go-scalingo"
)

var (
	scmIntegrationsListCommand = cli.Command{
		Name:     "scm-integrations",
		Category: "SCM Integrations",
		Usage:    "List your SCM integrations",
		Description: `List all the SCM integrations associated with your account:

	$ scalingo scm-integrations

	# See also commands 'scm-integrations-create', 'scm-integrations-delete', 'scm-integrations-import-keys'`,

		Action: func(c *cli.Context) {
			err := scm_integrations.List()
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "scm-integrations")
		},
	}

	scmIntegrationsCreateCommand = cli.Command{
		Name:     "scm-integrations-create",
		Category: "SCM Integrations",
		Flags: []cli.Flag{appFlag,
			cli.StringFlag{Name: "url", Usage: "URL of the SCM integration"},
			cli.StringFlag{Name: "token", Usage: "Token of the SCM integration"},
		},
		Usage: "Link your Scalingo account with your account on a SCM tool",
		Description: `Link your Scalingo account with your account on a SCM tool:

	For github.com:
	$ scalingo scm-integrations-create github

	For gitlab.com:
	$ scalingo scm-integrations-create gitlab

	For GitHub Enterprise:
	$ scalingo scm-integrations-create github-enterprise --url https://ghe.example.com --token personal-access-token

	For GitLab Self-hosted:
	$ scalingo scm-integrations-create gitlab-self-hosted --url https://gitlab.example.com --token personal-access-token

	# See also commands 'scm-integrations', 'scm-integrations-delete', 'scm-integrations-import-keys'`,

		Action: func(c *cli.Context) {
			if c.NArg() != 1 {
				_ = cli.ShowCommandHelp(c, "scm-integrations-create")
				return
			}

			link := c.String("url")
			token := c.String("token")
			scmType := scalingo.SCMType(c.Args()[0])

			switch scmType {
			case scalingo.SCMGithubType, scalingo.SCMGitlabType:
				break
			case scalingo.SCMGithubEnterpriseType, scalingo.SCMGitlabSelfHostedType:
				if link == "" || token == "" {
					errorQuit(errgo.New("both --url and --token must be set"))
				}

				u, err := url.Parse(link)
				if err != nil || u.Scheme == "" || u.Host == "" {
					errorQuit(errgo.Newf("'%s' is not a valid URL", link))
				}
			default:
				errorQuit(errgo.New(
					"Unknown SCM integration, available SCM integrations: github, github-enterprise, gitlab, gitlab-self-hosted",
				))
			}

			args := scm_integrations.CreateArgs{
				SCMType: scmType,
				URL:     link,
				Token:   token,
			}

			err := scm_integrations.Create(args)
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "scm-integrations-create")
		},
	}

	scmIntegrationsDeleteCommand = cli.Command{
		Name:     "scm-integrations-delete",
		Category: "SCM Integrations",
		Usage:    "Unlink your Scalingo account and your account on a SCM tool",
		Description: `Unlink your Scalingo account and your account on a SCM tool:

	$ scalingo scm-integrations-delete integration-type
	OR
	$ scalingo scm-integrations-delete integration-uuid

	Examples:
	$ scalingo scm-integrations-delete github-enterprise
	$ scalingo scm-integrations-delete gitlab

	# See also commands 'scm-integrations', 'scm-integrations-create', 'scm-integrations-import-keys'`,

		Action: func(c *cli.Context) {
			if c.NArg() != 1 {
				_ = cli.ShowCommandHelp(c, "scm-integrations-delete")
				return
			}

			err := scm_integrations.Delete(c.Args()[0])
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "scm-integrations-delete")
		},
	}

	scmIntegrationsImportKeysCommand = cli.Command{
		Name:     "scm-integrations-import-keys",
		Category: "SCM Integrations",
		Usage:    "Import public SSH keys from SCM account",
		Description: `Import public SSH keys from SCM account:

	$ scalingo scm-integrations-import-keys integration-type
	OR
	$ scalingo scm-integrations-import-keys integration-uuid

	Examples:
	$ scalingo scm-integrations-import-keys github
	$ scalingo scm-integrations-import-keys gitlab-self-hosted

	# See also commands 'scm-integrations', 'scm-integrations-create', 'scm-integrations-delete'`,

		Action: func(c *cli.Context) {
			if c.NArg() != 1 {
				_ = cli.ShowCommandHelp(c, "scm-integrations-import-keys")
				return
			}

			err := scm_integrations.ImportKeys(c.Args()[0])
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "scm-integrations-import-keys")
		},
	}
)
