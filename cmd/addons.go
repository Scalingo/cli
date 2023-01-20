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
		Description: CommandDescription{
			Description: "List all addons used by your app",
			Examples:    []string{"scalingo --app my-app addons"},
			SeeAlso:     []string{"addons-add", "addons-remove"},
		}.Render(),
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
		ArgsUsage: "addon-id plan",
		Description: CommandDescription{
			Description: "Provision an add-on for your application",
			Examples:    []string{"scalingo --app my-app addons-add mongodb mongo-starter-512"},
			SeeAlso:     []string{"addons-list", "addons-plans"},
		}.Render(),

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
		Name:      "addons-remove",
		Category:  "Addons",
		Flags:     []cli.Flag{&appFlag},
		Usage:     "Remove an existing addon from your app",
		ArgsUsage: "addon-id",
		Description: CommandDescription{
			Description: "Remove an existing addon from your app",
			Examples:    []string{"scalingo --app my-app addons-remove mongodb"},
			SeeAlso:     []string{"addons", "addons-add"},
		}.Render(),

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
		Name:      "addons-upgrade",
		Category:  "Addons",
		Flags:     []cli.Flag{&appFlag},
		Usage:     "Upgrade or downgrade an add-on attached to your app",
		ArgsUsage: "addon-id plan",
		Description: CommandDescription{
			Description: "Upgrade an addon attached to your app",
			Examples:    []string{"scalingo --app my-app addons-upgrade mongodb mongo-starter-256"},
			SeeAlso:     []string{"addons-plans", "addons-remove"},
		}.Render(),

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
		Name:      "addons-info",
		Category:  "Addons",
		Usage:     "Display information about an add-on attached to your app",
		Flags:     []cli.Flag{&appFlag},
		ArgsUsage: "addon-id",
		Description: CommandDescription{
			Description: "Display information about an add-on attached to your app",
			Examples:    []string{"scalingo --app my-app addons-info mongodb"},
			SeeAlso:     []string{"addons", "addons-upgrade"},
		}.Render(),
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
