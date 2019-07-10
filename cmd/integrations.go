package cmd

import (
	"github.com/urfave/cli"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/integrations"
)

var (
	IntegrationsListCommand = cli.Command{
		Name:     "integrations",
		Category: "Integrations",
		Usage:    "List your external integrations",
		Description: `List all the external integrations associated with your account:

	$ scalingo integrations

	# See also commands 'integrations-create', 'integrations-destroy', 'integrations-import-keys'`,

		Action: func(c *cli.Context) {
			err := integrations.List()
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "integrations")
		},
	}

	IntegrationsCreateCommand = cli.Command{
		Name:     "integrations-create",
		Category: "Integrations",
		Flags: []cli.Flag{appFlag,
			cli.StringFlag{Name: "url", Usage: "URL of the integration", Value: "<url>", EnvVar: ""},
			cli.StringFlag{Name: "token", Usage: "Token of the integration", Value: "<token>", EnvVar: ""},
		},
		Usage: "Create a link between external integration and your account",
		Description: `Create a link between external integration and your account:

	$ scalingo integrations-create --type integration-type --url integration-url --token integration-token

	# See also commands 'integrations', 'integrations-destroy', 'integrations-import-keys'`,

		Action: func(c *cli.Context) {
			var err error
			if len(c.Args()) == 1 {
				url := c.String("url")
				if url == "<url>" {
					url = ""
				}

				token := c.String("token")
				if token == "<token>" {
					token = ""
				}

				err = integrations.Create(c.Args()[0], url, token)
			} else {
				_ = cli.ShowCommandHelp(c, "integrations-create")
			}

			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "integrations-create")
		},
	}

	IntegrationsDestroyCommand = cli.Command{
		Name:     "integrations-destroy",
		Category: "Integrations",
		Usage:    "Destroy a link between external integration and your account",
		Description: `Destroy a link between external integration and your account:

	$ scalingo integrations-destroy integration-type OR scalingo integrations-import-keys integration-uuid

	# See also commands 'integrations', 'integrations-create', 'integrations-import-keys'`,

		Action: func(c *cli.Context) {
			var err error
			if len(c.Args()) == 1 {
				err = integrations.Destroy(c.Args()[0])
			} else {
				_ = cli.ShowCommandHelp(c, "integrations-destroy")
			}

			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "integrations-destroy")
		},
	}

	IntegrationsImportKeysCommand = cli.Command{
		Name:     "integrations-import-keys",
		Category: "Integrations",
		Usage:    "Import public SSH keys from external integration",
		Description: `Import public SSH keys from external integration:

	$ scalingo integrations-import-keys integration-type OR scalingo integrations-import-keys integration-uuid

	# See also commands 'integrations', 'integrations-create', 'integrations-destroy'`,

		Action: func(c *cli.Context) {
			var err error
			if len(c.Args()) == 1 {
				err = integrations.ImportKeys(c.Args()[0])
			} else {
				_ = cli.ShowCommandHelp(c, "integrations-import-keys")
			}

			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "integrations-import-keys")
		},
	}
)
