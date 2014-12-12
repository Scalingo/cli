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
)
