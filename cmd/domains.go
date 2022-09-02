package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/domains"
)

var (
	DomainsListCommand = cli.Command{
		Name:     "domains",
		Category: "Custom Domains",
		Flags:    []cli.Flag{&appFlag},
		Usage:    "List the domains of an application",
		Description: `List all the custom domains of an application:

    $ scalingo --app my-app domains

    # See also commands 'domains-add' and 'domains-remove'`,

		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			var err error
			if c.Args().Len() == 0 {
				err = domains.List(c.Context, currentApp)
			} else {
				cli.ShowCommandHelp(c, "domains")
			}

			if err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "domains")
		},
	}

	DomainsAddCommand = cli.Command{
		Name:     "domains-add",
		Category: "Custom Domains",
		Usage:    "Add a custom domain to an application",
		Flags: []cli.Flag{&appFlag,
			&cli.StringFlag{Name: "cert", Usage: "SSL Signed Certificate", Value: "domain.crt"},
			&cli.StringFlag{Name: "key", Usage: "SSL Keypair", Value: "domain.key"},
		},
		Description: `Add a custom domain to an application:

    $ scalingo -a myapp domains-add example.com

    # See also commands 'domains' and 'domains-remove'`,

		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
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
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "domains-add")
		},
	}

	DomainsRemoveCommand = cli.Command{
		Name:     "domains-remove",
		Category: "Custom Domains",
		Flags:    []cli.Flag{&appFlag},
		Usage:    "Remove a custom domain from an application",
		Description: `Remove a custom domain from an application:

    $ scalingo -a myapp domains-remove example.com

    # See also commands 'domains' and 'domains-add'`,

		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			var err error
			if c.Args().Len() == 1 {
				err = domains.Remove(c.Context, currentApp, c.Args().First())
			} else {
				cli.ShowCommandHelp(c, "domains-remove")
			}

			if err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "domains-remove")
			autocomplete.DomainsRemoveAutoComplete(c)
		},
	}

	DomainsSSLCommand = cli.Command{
		Name:     "domains-ssl",
		Category: "Custom Domains",
		Usage:    "Enable or disable SSL for your custom domains",
		Flags: []cli.Flag{&appFlag,
			&cli.StringFlag{Name: "cert", Usage: "SSL Signed Certificate", Value: "domain.crt"},
			&cli.StringFlag{Name: "key", Usage: "SSL Keypair", Value: "domain.key"},
		},
		Description: `Enable or disable SSL for your custom domains:

		$ scalingo -a myapp domains-ssl example.com --cert <certificate.crt> --key <keyfile.key>

		$ scalingo -a myapp domains-ssl example.com disable

		# See also commands 'domains' and 'domains-add'`,

		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			var err error
			if c.Args().Len() == 2 && c.Args().Slice()[1] == "disable" {
				err = domains.DisableSSL(c.Context, currentApp, c.Args().First())
			} else if c.Args().Len() == 1 {
				err = domains.EnableSSL(c.Context, currentApp, c.Args().First(), c.String("cert"), c.String("key"))
			} else {
				cli.ShowCommandHelp(c, "domains-ssl")
			}
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "domains-ssl")
		},
	}

	setCanonicalDomainCommand = cli.Command{
		Name:     "set-canonical-domain",
		Category: "App Management",
		Flags:    []cli.Flag{&appFlag},
		Usage:    "Set a canonical domain.",
		Description: `After defining multiple domains on an application, one can need to redirect all requests towards its application to a specific domain. This domain is called the canonical domain. This command sets the canonical domain:

    $ scalingo -a myapp set-canonical-domain example.com

    # See also commands 'domains', 'domains-add' and 'unset-canonical-domain'`,

		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() != 1 {
				cli.ShowCommandHelp(c, "set-canonical-domain")
				return nil
			}

			err := domains.SetCanonical(c.Context, currentApp, c.Args().First())
			if err != nil {
				errorQuit(err)
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
		Description: `Unset the canonical domain of this app:

    $ scalingo -a myapp unset-canonical-domain

    # See also commands 'domains', 'domains-add' and 'set-canonical-domain'`,

		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() != 0 {
				cli.ShowCommandHelp(c, "unset-canonical-domain")
				return nil
			}

			err := domains.UnsetCanonical(c.Context, currentApp)
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "unset-canonical-domain")
		},
	}
)
