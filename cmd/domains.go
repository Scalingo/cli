package cmd

import (
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/domains"
	"github.com/codegangsta/cli"
)

var (
	DomainsListCommand = cli.Command{
		Name:  "domains",
		Usage: "List the domains of an application",
		Description: `List all the custom domains of an application:

    $ scalingo -a myapp domains

    # See also commands 'domains-add' and 'domains-remove'`,

		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c.GlobalString("app"))
			var err error
			if len(c.Args()) == 0 {
				err = domains.List(currentApp)
			} else {
				cli.ShowCommandHelp(c, "domains")
			}

			if err != nil {
				errorQuit(err)
			}
		},
	}

	DomainsAddCommand = cli.Command{
		Name:  "domains-add",
		Usage: "Add a custom domain to an application",
		Description: `Add a custom domain to an application:

    $ scalingo -a myapp domains-add example.com

    # See also commands 'domains' and 'domains-remove'`,

		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c.GlobalString("app"))
			var err error
			if len(c.Args()) == 1 {
				err = domains.Add(currentApp, c.Args()[0])
			} else {
				cli.ShowCommandHelp(c, "domains-add")
			}

			if err != nil {
				errorQuit(err)
			}
		},
	}

	DomainsRemoveCommand = cli.Command{
		Name:  "domains-remove",
		Usage: "Remove a custom domain from an application",
		Description: `Remove a custom domain from an application:

    $ scalingo -a myapp domains-remove example.com

    # See also commands 'domains' and 'domains-add'`,

		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c.GlobalString("app"))
			var err error
			if len(c.Args()) == 1 {
				err = domains.Remove(currentApp, c.Args()[0])
			} else {
				cli.ShowCommandHelp(c, "domains-remove")
			}

			if err != nil {
				errorQuit(err)
			}
		},
	}

	DomainsSSLCommand = cli.Command{
		Name:  "domains-ssl",
		Usage: "Enable or disable SSL for your custom domains",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "cert", Usage: "SSL Signed Certificate", Value: "domain.crt", EnvVar: ""},
			cli.StringFlag{Name: "key", Usage: "SSL Keypair", Value: "domain.key", EnvVar: ""},
		},
		Description: `Enable or disable SSL for your custom domains:
		
		$ scalingo -a myapp domains-ssl example.com --cert <certificate.crt> --key <keyfile.key>

		$ scalingo -a myapp domains-ssl example.com disable
		
		# See also commands 'domains' and 'domains-add'`,

		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c.GlobalString("app"))
			var err error
			if len(c.Args()) == 2 && c.Args()[1] == "disable" {
				err = domains.DisableSSL(currentApp, c.Args()[0])
			} else if len(c.Args()) == 1 {
				err = domains.EnableSSL(currentApp, c.Args()[0], c.String("cert"), c.String("key"))
			} else {
				cli.ShowCommandHelp(c, "domains-ssl")
			}
			if err != nil {
				errorQuit(err)
			}
		},
	}
)
