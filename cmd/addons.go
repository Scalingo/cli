package cmd

import (
	"github.com/urfave/cli"

	"github.com/Scalingo/cli/addons"
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/cmd/autocomplete"
)

var (
	AddonsListCommand = cli.Command{
		Name:     "addons",
		Category: "Addons",
		Usage:    "List used add-ons",
		Flags:    []cli.Flag{appFlag},
		Description: ` List all addons used by your app:
    $ scalingo -a myapp addons

		# See also 'addons-add' and 'addons-remove'
`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
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
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "addons")
		},
	}
	AddonsAddCommand = cli.Command{
		Name:     "addons-add",
		Category: "Addons",
		Flags:    []cli.Flag{appFlag},
		Usage:    "Provision an add-on for your application",
		Description: ` Provision an add-on for your application:
    $ scalingo -a myapp addons-add <addon-name> <plan>

		# See also 'addons-list' and 'addons-plans'
`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
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
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "addons-add")
			autocomplete.AddonsAddAutoComplete(c)
		},
	}
	AddonsRemoveCommand = cli.Command{
		Name:     "addons-remove",
		Category: "Addons",
		Flags:    []cli.Flag{appFlag},
		Usage:    "Remove an existing addon from your app",
		Description: ` Remove an existing addon from your app:
    $ scalingo -a myapp addons-remove <addon-id>

		# See also 'addons' and 'addons-add'
`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
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
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "addons-remove")
			autocomplete.AddonsRemoveAutoComplete(c)
		},
	}
	AddonsUpgradeCommand = cli.Command{
		Name:     "addons-upgrade",
		Category: "Addons",
		Flags:    []cli.Flag{appFlag},
		Usage:    "Upgrade or downgrade an add-on attached to your app",
		Description: ` Upgrade an addon attached to your app:
    $ scalingo -a myapp addons-upgrade <addon-id> <plan>

		# See also 'addons-plans' and 'addons-remove'
`,
		Action: func(c *cli.Context) {
			currentApp := appdetect.CurrentApp(c)
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
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "addons-upgrade")
			autocomplete.AddonsUpgradeAutoComplete(c)
		},
	}
	AddonsInfoCommand = cli.Command{
		Name:     "addons-info",
		Category: "Addons",
		Usage:    "Display infos about add-ons",
		Flags:    []cli.Flag{appFlag},
		Description: ` Display infos about add-ons:
    $ scalingo --app my-app addons-info --addon <addon-id>

		# See also 'addons' and 'addons-add'
`,
		Action: func(c *cli.Context) {
			if len(c.Args()) != 0 {
				cli.ShowCommandHelp(c, "addons-info")
				return
			}

			currentApp := appdetect.CurrentApp(c)
			currentAddon := addonName(c)

			err := addons.Info(currentApp, currentAddon)
			if err != nil {
				errorQuit(err)
			}
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "addons-info")
		},
	}
)
