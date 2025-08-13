package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/addons"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-utils/errors/v2"
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
		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() > 0 {
				return cli.ShowCommandHelp(ctx, c, "addons")
			}

			err := addons.List(ctx, currentApp)
			if err != nil {
				errorQuit(ctx, err)
			}

			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "addons")
		},
	}
	addonsAddCommand = cli.Command{
		Name:      "addons-add",
		Category:  "Addons",
		Flags:     []cli.Flag{&appFlag},
		Usage:     "Provision an add-on for your application",
		ArgsUsage: "addon-name plan",
		Description: CommandDescription{
			Description: "Provision an add-on for your application",
			Examples:    []string{"scalingo --app my-app addons-add postgresql postgresql-starter-1024"},
			SeeAlso:     []string{"addons-list", "addons-plans"},
		}.Render(),

		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() != 2 {
				return cli.ShowCommandHelp(ctx, c, "addons-add")
			}

			utils.CheckForConsent(ctx, currentApp)

			err := addons.Provision(ctx, currentApp, c.Args().First(), c.Args().Slice()[1])
			if err != nil {
				errorQuit(ctx, err)
			}

			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "addons-add")
			_ = autocomplete.AddonsAddAutoComplete(ctx)
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
			Examples:    []string{"scalingo --app my-app addons-remove addon_uuid"},
			SeeAlso:     []string{"addons", "addons-add"},
		}.Render(),

		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() != 1 {
				return cli.ShowCommandHelp(ctx, c, "addons-remove")
			}

			utils.CheckForConsent(ctx, currentApp)

			addonUUID, err := utils.GetAddonUUIDFromType(ctx, currentApp, c.Args().First())
			if err != nil {
				errorQuit(ctx, err)
			}

			err = addons.Destroy(ctx, currentApp, addonUUID)
			if err != nil {
				errorQuit(ctx, err)
			}

			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "addons-remove")
			_ = autocomplete.AddonsRemoveAutoComplete(ctx, c)
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
			Examples:    []string{"scalingo --app my-app addons-upgrade addon_uuid mongo-starter-256"},
			SeeAlso:     []string{"addons-plans", "addons-remove"},
		}.Render(),

		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() != 2 {
				return cli.ShowCommandHelp(ctx, c, "addons-upgrade")
			}

			addonUUID, err := utils.GetAddonUUIDFromType(ctx, currentApp, c.Args().First())
			if err != nil {
				errorQuit(ctx, err)
			}

			err = addons.Upgrade(ctx, currentApp, addonUUID, c.Args().Slice()[1])
			if err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "addons-upgrade")
			_ = autocomplete.AddonsUpgradeAutoComplete(ctx, c)
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
			Examples:    []string{"scalingo --app my-app addons-info addon_uuid"},
			SeeAlso:     []string{"addons", "addons-upgrade"},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() != 1 {
				return cli.ShowCommandHelp(ctx, c, "addons-info")
			}

			currentApp := detect.CurrentApp(c)
			addonName := c.Args().First()

			addonUUID, err := utils.GetAddonUUIDFromType(ctx, currentApp, addonName)
			if err != nil {
				errorQuit(ctx, err)
			}

			err = addons.Info(ctx, currentApp, addonUUID)
			if err != nil {
				errorQuit(ctx, err)
			}

			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "addons-info")
		},
	}
	addonsConfigCommand = cli.Command{
		Name:     "addons-config",
		Category: "Addons",
		Usage:    "Configure an add-on attached to your app",
		Flags: []cli.Flag{&appFlag,
			&cli.StringFlag{Name: "maintenance-window-hour", Usage: "Configure add-on maintenance window starting hour (in your local timezone)"},
			&cli.StringFlag{Name: "maintenance-window-day", Usage: "Configure add-on maintenance window day"},
		},
		ArgsUsage: "addon-id",
		Description: CommandDescription{
			Description: "Configure an add-on attached to your app",
			Examples:    []string{"scalingo --app my-app addons-config --maintenance-window-day tuesday --maintenance-window-hour 2 addon_uuid"},
			SeeAlso:     []string{"addons", "addons-info"},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() != 1 {
				return cli.ShowCommandHelp(ctx, c, "addons-config")
			}

			currentApp := detect.CurrentApp(c)
			currentAddon := c.Args().First()
			config := addons.UpdateAddonConfigOpts{}

			if c.IsSet("maintenance-window-day") {
				config.MaintenanceWindowDay = utils.StringPtr(c.String("maintenance-window-day"))
			}

			if c.IsSet("maintenance-window-hour") {
				config.MaintenanceWindowHour = utils.IntPtr(c.Int("maintenance-window-hour"))
			}

			if config.MaintenanceWindowHour == nil && config.MaintenanceWindowDay == nil {
				errorQuitWithHelpMessage(ctx, errors.New(ctx, "cannot update your addon without a specified option"), c, "addons-config")
			}

			err := addons.UpdateConfig(ctx, currentApp, currentAddon, config)
			if err != nil {
				errorQuit(ctx, err)
			}

			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			err := autocomplete.CmdFlagsAutoComplete(c, "addons-config")
			if err != nil {
				errorQuit(ctx, err)
			}
		},
	}
)
