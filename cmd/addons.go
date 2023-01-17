package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/cli/addons"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
)

var (
	addonsListCommand = cli.Command{
		Name:     "addons",
		Category: "Addons",
		Usage:    "List used add-ons",
		Flags:    []cli.Flag{&appFlag},
		Description: ` List all addons used by your app:
    $ scalingo -a myapp addons

		# See also 'addons-add' and 'addons-remove'
`,
		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() > 0 {
				return cli.ShowCommandHelp(c, "addons")
			}

			err := addons.List(c.Context, currentApp)
			if err != nil {
				errorQuit(err)
			}

			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "addons")
		},
	}
	addonsAddCommand = cli.Command{
		Name:      "addons-add",
		Category:  "Addons",
		Flags:     []cli.Flag{&appFlag},
		Usage:     "Provision an add-on for your application",
		ArgsUsage: "addon-name plan",
		Description: `Provision an add-on for your application:

Usage:
  $ scalingo -a myapp addons-add <addon-name> <plan>

# See also 'addons-list' and 'addons-plans'`,

		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() != 2 {
				return cli.ShowCommandHelp(c, "addons-add")
			}

			err := addons.Provision(c.Context, currentApp, c.Args().First(), c.Args().Slice()[1])
			if err != nil {
				errorQuit(err)
			}

			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "addons-add")
			autocomplete.AddonsAddAutoComplete(c)
		},
	}
	addonsRemoveCommand = cli.Command{
		Name:     "addons-remove",
		Category: "Addons",
		Flags:    []cli.Flag{&appFlag},
		Usage:    "Remove an existing addon from your app",
		Description: ` Remove an existing addon from your app:
    $ scalingo -a myapp addons-remove <addon-id>

		# See also 'addons' and 'addons-add'
`,
		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() != 1 {
				return cli.ShowCommandHelp(c, "addons-remove")
			}

			err := addons.Destroy(c.Context, currentApp, c.Args().First())
			if err != nil {
				errorQuit(err)
			}

			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "addons-remove")
			autocomplete.AddonsRemoveAutoComplete(c)
		},
	}
	addonsUpgradeCommand = cli.Command{
		Name:     "addons-upgrade",
		Category: "Addons",
		Flags:    []cli.Flag{&appFlag},
		Usage:    "Upgrade or downgrade an add-on attached to your app",
		Description: ` Upgrade an addon attached to your app:
    $ scalingo -a myapp addons-upgrade <addon-id> <plan>

		# See also 'addons-plans' and 'addons-remove'
`,
		Action: func(c *cli.Context) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() != 2 {
				return cli.ShowCommandHelp(c, "addons-upgrade")
			}

			err := addons.Upgrade(c.Context, currentApp, c.Args().First(), c.Args().Slice()[1])
			if err != nil {
				errorQuit(err)
			}
			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "addons-upgrade")
			autocomplete.AddonsUpgradeAutoComplete(c)
		},
	}
	addonsInfoCommand = cli.Command{
		Name:     "addons-info",
		Category: "Addons",
		Usage:    "Display information about an add-on attached to your app",
		Flags:    []cli.Flag{&appFlag},
		Description: `Display information about an add-on attached to your app:

$ scalingo --app my-app addons-info <addon-id>

# See also 'addons' and 'addons-upgrade'
`,
		Action: func(c *cli.Context) error {
			if c.Args().Len() != 1 {
				return cli.ShowCommandHelp(c, "addons-info")
			}

			currentApp := detect.CurrentApp(c)
			currentAddon := c.Args().First()

			err := addons.Info(c.Context, currentApp, currentAddon)
			if err != nil {
				errorQuit(err)
			}

			return nil
		},
		BashComplete: func(c *cli.Context) {
			autocomplete.CmdFlagsAutoComplete(c, "addons-info")
		},
	}
)
