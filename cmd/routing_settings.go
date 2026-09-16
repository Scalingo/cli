package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-utils/errors/v3"
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
			currentApp := detect.CurrentApp(ctx, c)
			if c.Args().Len() > 1 {
				_ = cli.ShowCommandHelp(ctx, c, "force-https")
				return nil
			}

			utils.CheckForConsent(ctx, currentApp, utils.ConsentTypeContainers)

			enable := !c.IsSet("disable")

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
			currentApp := detect.CurrentApp(ctx, c)
			if c.Args().Len() > 1 {
				_ = cli.ShowCommandHelp(ctx, c, "sticky-session")
				return nil
			}

			utils.CheckForConsent(ctx, currentApp, utils.ConsentTypeContainers)

			enable := !c.IsSet("disable")

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
			currentApp := detect.CurrentApp(ctx, c)
			if c.Args().Len() > 1 {
				_ = cli.ShowCommandHelp(ctx, c, "router-logs")
				return nil
			}

			utils.CheckForConsent(ctx, currentApp, utils.ConsentTypeContainers)

			enable := c.IsSet("enable")

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

	appFirewallRulesCommand = cli.Command{
		Name:     "app-firewall-rules",
		Category: "App Management",
		Usage:    "List IPv4 CIDR firewall rules for your application",
		Flags: []cli.Flag{
			&appFlag,
		},
		Description: CommandDescription{
			Description: "List app-level IPv4 CIDR firewall rules.",
			Examples: []string{
				"scalingo --app my-app app-firewall-rules",
			},
			SeeAlso: []string{"app-firewall-rule-add", "app-firewall-rule-remove"},
		}.Render(),

		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(ctx, c)
			if c.Args().Len() > 0 {
				_ = cli.ShowCommandHelp(ctx, c, "app-firewall-rules")
				return nil
			}

			err := apps.AppFirewallRulesList(ctx, currentApp)
			if err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "app-firewall-rules")
		},
	}

	appFirewallRuleAddCommand = cli.Command{
		Name:      "app-firewall-rule-add",
		Category:  "App Management",
		Usage:     "Add an IPv4 CIDR firewall rule for your application",
		ArgsUsage: "cidr",
		Flags: []cli.Flag{
			&appFlag,
			&cli.StringFlag{Name: "cidr", Usage: "IPv4 CIDR to allow", Required: true},
			&cli.StringFlag{Name: "label", Usage: "Optional label attached to the rule"},
		},
		Description: CommandDescription{
			Description: "Add an app-level IPv4 CIDR firewall rule.",
			Examples: []string{
				"scalingo --app my-app app-firewall-rule-add 203.0.113.42/32 --label office",
				"scalingo --app my-app app-firewall-rule-add --cidr 10.0.0.0/24",
			},
			SeeAlso: []string{"app-firewall-rules", "app-firewall-rule-remove"},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(ctx, c)
			cidr := c.String("cidr")
			if cidr == "" && c.Args().Len() == 1 {
				cidr = c.Args().First()
			}
			if cidr == "" || c.Args().Len() > 1 {
				errorQuitWithHelpMessage(ctx, errors.New(ctx, "a CIDR argument or --cidr is required"), c, "app-firewall-rule-add")
				return nil
			}

			err := apps.AppFirewallRuleAdd(ctx, currentApp, cidr, c.String("label"))
			if err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "app-firewall-rule-add")
		},
	}

	appFirewallRuleRemoveCommand = cli.Command{
		Name:      "app-firewall-rule-remove",
		Category:  "App Management",
		Usage:     "Remove an IPv4 CIDR firewall rule from your application",
		ArgsUsage: "rule-id",
		Flags: []cli.Flag{
			&appFlag,
		},
		Description: CommandDescription{
			Description: "Remove an app-level IPv4 CIDR firewall rule. The rule ID is shown by app-firewall-rules.",
			Examples: []string{
				"scalingo --app my-app app-firewall-rule-remove afwr-123",
			},
			SeeAlso: []string{"app-firewall-rules", "app-firewall-rule-add"},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(ctx, c)
			if c.Args().Len() != 1 {
				_ = cli.ShowCommandHelp(ctx, c, "app-firewall-rule-remove")
				return nil
			}

			err := apps.AppFirewallRuleRemove(ctx, currentApp, c.Args().First())
			if err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "app-firewall-rule-remove")
		},
	}
)
