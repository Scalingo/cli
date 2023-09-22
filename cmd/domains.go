package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/domains"
	"github.com/Scalingo/cli/utils"
)

var (
	DomainsListCommand = cli.Command{
		Name:     "domains",
		Category: "Custom Domains",
		Flags:    []cli.Flag{&appFlag},
		Usage:    "List the domains of an application",
		Description: CommandDescription{
			Description: "List all the custom domains of an application",
			Examples:    []string{"scalingo --app my-app domains"},
			SeeAlso:     []string{"domains-add", "domains-remove"},
		}.Render(),

		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			var err error
			if c.Args().Len() == 0 {
				err = domains.List(c.Context, currentApp)
			} else {
				cli.ShowCommandHelp(c, "domains")
			}

			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "domains")
		},
	}

	DomainsAddCommand = cli.Command{
		Name:      "domains-add",
		Category:  "Custom Domains",
		Usage:     "Add a custom domain to an application",
		ArgsUsage: "domain",
		Flags: []cli.Flag{&appFlag,
			&cli.StringFlag{Name: "cert", Usage: "SSL Signed Certificate", Value: "domain.crt"},
			&cli.StringFlag{Name: "key", Usage: "SSL Keypair", Value: "domain.key"},
		},
		Description: CommandDescription{
			Description: "Add a custom domain to an application",
			Examples:    []string{"scalingo --app my-app domains-add example.com"},
			SeeAlso:     []string{"domains", "domains-remove"},
		}.Render(),

		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeContainers)
			var err error
			if c.Args().Len() == 1 {
				cert := c.String("cert")
				if cert == "domain.crt" {
					cert = ""
				}
				key := c.String("key")
				if key == "domain.key" {
					key = ""
				}
				err = domains.Add(c.Context, currentApp, c.Args().First(), cert, key)
			} else {
				cli.ShowCommandHelp(c, "domains-add")
			}

			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "domains-add")
		},
	}

	DomainsRemoveCommand = cli.Command{
		Name:      "domains-remove",
		Category:  "Custom Domains",
		Flags:     []cli.Flag{&appFlag},
		Usage:     "Remove a custom domain from an application",
		ArgsUsage: "domain",
		Description: CommandDescription{
			Description: "Remove a custom domain from an application",
			Examples:    []string{"scalingo --app my-app domains-remove example.com"},
			SeeAlso:     []string{"domains", "domains-add"},
		}.Render(),

		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeContainers)
			var err error
			if c.Args().Len() == 1 {
				err = domains.Remove(c.Context, currentApp, c.Args().First())
			} else {
				cli.ShowCommandHelp(c, "domains-remove")
			}

			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "domains-remove")
			autocomplete.DomainsRemoveAutoComplete(c)
		},
	}

	// TODO: Split the two operations (enable/disable) into two subcommands
	DomainsSSLCommand = cli.Command{
		Name:      "domains-ssl",
		Category:  "Custom Domains",
		Usage:     "Enable or disable SSL for your custom domains",
		ArgsUsage: "domain [disable]",
		Flags: []cli.Flag{&appFlag,
			&cli.StringFlag{Name: "cert", Usage: "SSL Signed Certificate", Value: "domain.crt"},
			&cli.StringFlag{Name: "key", Usage: "SSL Keypair", Value: "domain.key"},
		},
		Description: CommandDescription{
			Description: "Enable or disable SSL for your custom domains",
			Examples: []string{
				"scalingo --app my-app domains-ssl --cert certificate.crt --key keyfile.key example.com",
				"scalingo --app my-app domains-ssl example.com disable",
			},
			SeeAlso: []string{"domains", "domains-add"},
		}.Render(),

		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeContainers)
			var err error
			if c.Args().Len() == 2 && c.Args().Slice()[1] == "disable" {
				err = domains.DisableSSL(c.Context, currentApp, c.Args().First())
			} else if c.Args().Len() == 1 {
				err = domains.EnableSSL(c.Context, currentApp, c.Args().First(), c.String("cert"), c.String("key"))
			} else {
				cli.ShowCommandHelp(c, "domains-ssl")
			}
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "domains-ssl")
		},
	}

	setCanonicalDomainCommand = cli.Command{
		Name:      "set-canonical-domain",
		Category:  "App Management",
		Flags:     []cli.Flag{&appFlag},
		Usage:     "Set a canonical domain.",
		ArgsUsage: "domain",
		Description: CommandDescription{
			Description: `After defining multiple domains on an application, one can need to redirect all requests towards its application to a specific domain.
This domain is called the canonical domain. This command sets the canonical domain`,
			Examples: []string{"scalingo -a myapp set-canonical-domain example.com"},
			SeeAlso:  []string{"domains", "domains-add", "unset-canonical-domain"},
		}.Render(),

		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() != 1 {
				cli.ShowCommandHelp(c, "set-canonical-domain")
				return nil
			}
			utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeContainers)

			err := domains.SetCanonical(c.Context, currentApp, c.Args().First())
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "set-canonical-domain")
		},
	}

	unsetCanonicalDomainCommand = cli.Command{
		Name:     "unset-canonical-domain",
		Category: "App Management",
		Flags:    []cli.Flag{&appFlag},
		Usage:    "Unset a canonical domain.",
		Description: CommandDescription{
			Description: "Unset the canonical domain of this app",
			Examples:    []string{"scalingo --app my-app unset-canonical-domain"},
			SeeAlso:     []string{"domains", "domains-add", "set-canonical-domain"},
		}.Render(),

		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() != 0 {
				cli.ShowCommandHelp(c, "unset-canonical-domain")
				return nil
			}

			utils.CheckForConsent(c.Context, currentApp, utils.ConsentTypeContainers)

			err := domains.UnsetCanonical(c.Context, currentApp)
			if err != nil {
				errorQuit(c.Context, err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "unset-canonical-domain")
		},
	}
)
