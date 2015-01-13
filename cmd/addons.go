package cmd

import (
	"github.com/Scalingo/cli/addons"
	"github.com/Scalingo/cli/appdetect"
	"github.com/codegangsta/cli"
)

var (
	AddonsListCommand = cli.Command{
		Name:     "addons",
		Category: "Addons",
		Usage:    "List used add-ons",
		Description: ` List all addons used by your app:
    $ scalingo -a myapp addons

		# See also 'addons-add' and 'addons-remove'
`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c.GlobalString("app"))
			var err error
			if len(c.Args()) == 0 {
				err = addons.List(currentApp)
			} else {
				cli.ShowCommandHelp(c, "addons")
			}

			if err != nil {
				errorQuit(err)
			}
		},
	}
	AddonsAddCommand = cli.Command{
		Name:     "addons-add",
		Category: "Addons",
		Usage:    "Provision an add-on for your application",
		Description: ` Provision an add-on for your application:
    $ scalingo -a myapp addons-add <addon-name> <plan>

		# See also 'addons-list' and 'addons-plans'
`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c.GlobalString("app"))
			var err error
			if len(c.Args()) == 2 {
				err = addons.Provision(currentApp, c.Args()[0], c.Args()[1])
			} else {
				cli.ShowCommandHelp(c, "addons-add")
			}
			if err != nil {
				errorQuit(err)
			}
		},
	}
	AddonsRemoveCommand = cli.Command{
		Name:     "addons-remove",
		Category: "Addons",
		Usage:    "Remove an existing addon from your app",
		Description: ` Remove an existing addon from your app:
    $ scalingo -a myapp addons-remove <addon-id>

		# See also 'addons' and 'addons-remove'
`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c.GlobalString("app"))
			var err error
			if len(c.Args()) == 1 {
				err = addons.Destroy(currentApp, c.Args()[0])
			} else {
				cli.ShowCommandHelp(c, "addons-remove")
			}
			if err != nil {
				errorQuit(err)
			}
		},
	}
	AddonsUpgradeCommand = cli.Command{
		Name:     "addons-upgrade",
		Category: "Addons",
		Usage:    "Upgrade or downgrade an add-on attached to your app",
		Description: ` Upgrade an addon attached to your app:
    $ scalingo -a myapp addons-upgrade <addon-id> <plan>

		# See also 'addons-plans' and 'addons-remove'
`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c.GlobalString("app"))
			var err error
			if len(c.Args()) == 2 {
				err = addons.Upgrade(currentApp, c.Args()[0], c.Args()[1])
			} else {
				cli.ShowCommandHelp(c, "addons-upgrade")
			}
			if err != nil {
				errorQuit(err)
			}
		},
	}
)
