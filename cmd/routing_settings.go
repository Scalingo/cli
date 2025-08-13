package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/utils"
)

var (
	forceHTTPSCommand = cli.Command{
		Name:     "force-https",
		Category: "App Management",
		Usage:    "Enable/Disable automatic redirection of traffic to HTTPS for your application",
		Flags: []cli.Flag{
			&appFlag,
			&cli.BoolFlag{Name: "enable", Aliases: []string{"e"}, Usage: "Enable force HTTPS (default)"},
			&cli.BoolFlag{Name: "disable", Aliases: []string{"d"}, Usage: "Disable force HTTPS"},
		},
		Description: CommandDescription{
			Description: "When enabled, this feature will automatically redirect HTTP traffic to HTTPS for all domains associated with this application.",
			Examples:    []string{"scalingo --app my-app force-https --enable"},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() > 1 {
				_ = cli.ShowCommandHelp(ctx, c, "force-https")
				return nil
			}

			utils.CheckForConsent(ctx, currentApp, utils.ConsentTypeContainers)

			enable := true
			if c.IsSet("disable") {
				enable = false
			}

			err := apps.ForceHTTPS(ctx, currentApp, enable)
			if err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "force-https")
		},
	}

	stickySessionCommand = cli.Command{
		Name:     "sticky-session",
		Category: "App Management",
		Usage:    "Enable/Disable sticky sessions for your application",
		Flags: []cli.Flag{
			&appFlag,
			&cli.BoolFlag{Name: "enable", Aliases: []string{"e"}, Usage: "Enable sticky session (default)"},
			&cli.BoolFlag{Name: "disable", Aliases: []string{"d"}, Usage: "Disable sticky session"},
		},
		Description: CommandDescription{
			Description: "When enabled, application user sessions will be sticky: they will always return to the same container",
			Examples:    []string{"scalingo --app my-app sticky-session --enable"},
		}.Render(),

		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() > 1 {
				_ = cli.ShowCommandHelp(ctx, c, "sticky-session")
				return nil
			}

			utils.CheckForConsent(ctx, currentApp, utils.ConsentTypeContainers)

			enable := true
			if c.IsSet("disable") {
				enable = false
			}

			err := apps.StickySession(ctx, currentApp, enable)
			if err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "sticky-session")
		},
	}

	routerLogsCommand = cli.Command{
		Name:     "router-logs",
		Category: "App Management",
		Usage:    "Enable/disable router logs for your application",
		Flags: []cli.Flag{
			&appFlag,
			&cli.BoolFlag{Name: "enable", Aliases: []string{"e"}, Usage: "Enable router logs"},
			&cli.BoolFlag{Name: "disable", Aliases: []string{"d"}, Usage: "Disable router logs (default)"},
		},
		Description: CommandDescription{
			Description: "Enable/disable router logs for your application",
			Examples:    []string{"scalingo --app my-app router-logs --enable"},
		}.Render(),

		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() > 1 {
				_ = cli.ShowCommandHelp(ctx, c, "router-logs")
				return nil
			}

			utils.CheckForConsent(ctx, currentApp, utils.ConsentTypeContainers)

			enable := false
			if c.IsSet("enable") {
				enable = true
			}

			err := apps.RouterLogs(ctx, currentApp, enable)
			if err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "router-logs")
		},
	}
)
