package cmd

import (
	"net/url"

	"github.com/urfave/cli"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/scm_integrations"
)

var (
	scmIntegrationsListCommand = cli.Command{
		Name:     "scm-integrations",
		Category: "SCM Integrations",
		Usage:    "List your scm integrations",
		Description: `List all the scm integrations associated with your account:

	$ scalingo scm-integrations

	# See also commands 'scm-integrations-create', 'scm-integrations-destroy', 'scm-integrations-import-keys'`,

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
			cli.StringFlag{Name: "url", Usage: "URL of the scm integration"},
			cli.StringFlag{Name: "token", Usage: "Token of the scm integration"},
		},
		Usage: "Create a link between scm integration and your account",
		Description: `Create a link between scm integration and your account:

	For github.com:
	$ scalingo scm-integrations-create github

	For gitlab.com:
	$ scalingo scm-integrations-create gitlab

	For GitHub Enterprise:
	$ scalingo scm-integrations-create github-enterprise --url https://ghe.example.com --token personal-access-token

	For GitLab Self-hosted:
	$ scalingo scm-integrations-create gitlab-self-hosted --url https://gitlab.example.com --token personal-access-token

	# See also commands 'scm-integrations', 'scm-integrations-destroy', 'scm-integrations-import-keys'`,

		Action: func(c *cli.Context) {
			if c.NArg() != 1 {
				_ = cli.ShowCommandHelp(c, "scm-integrations-create")
				return
			}

			link := c.String("url")
			token := c.String("token")

			switch c.Args()[0] {
			case "github", "gitlab":
				break
			case "github-enterprise", "gitlab-self-hosted":
				if link == "" || token == "" {
					errorQuit(errgo.New("URL or Token is not set"))
				}

				u, err := url.Parse(link)
				if err != nil || u.Scheme == "" || u.Host == "" {
					errorQuit(errgo.New("URL is not a valid url"))
				}
			default:
				errorQuit(errgo.New(
					"Unknown scm integration, available scm integrations : github, github-enterprise, gitlab, gitlab-self-hosted",
				))
			}

			args := scm_integrations.CreateArgs{
				ScmType: c.Args()[0],
				Url:     link,
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

	scmIntegrationsDestroyCommand = cli.Command{
		Name:     "scm-integrations-destroy",
		Category: "SCM Integrations",
		Usage:    "Destroy a link between scm integration and your account",
		Description: `Destroy a link between scm integration and your account:

	$ scalingo scm-integrations-destroy integration-type
	OR
	$ scalingo scm-integrations-destroy integration-uuid

	Examples:
	$ scalingo scm-integrations-destroy github-enterprise
	$ scalingo scm-integrations-destroy gitlab

	# See also commands 'scm-integrations', 'scm-integrations-create', 'scm-integrations-import-keys'`,

		Action: func(c *cli.Context) {
			if c.NArg() != 1 {
				_ = cli.ShowCommandHelp(c, "scm-integrations-destroy")
				return
			}

			err := scm_integrations.Destroy(c.Args()[0])
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "scm-integrations-destroy")
		},
	}

	scmIntegrationsImportKeysCommand = cli.Command{
		Name:     "scm-integrations-import-keys",
		Category: "SCM Integrations",
		Usage:    "Import public SSH keys from scm integration",
		Description: `Import public SSH keys from scm integration:

	$ scalingo scm-integrations-import-keys integration-type
	OR
	$ scalingo scm-integrations-import-keys integration-uuid

	Examples:
	$ scalingo scm-integrations-import-keys github
	$ scalingo scm-integrations-import-keys gitlab-self-hosted

	# See also commands 'scm-integrations', 'scm-integrations-create', 'scm-integrations-destroy'`,

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
