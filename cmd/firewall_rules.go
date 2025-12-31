package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/dbng"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v8"
)

var (
	databaseFirewallRulesListCommand = cli.Command{
		Name:      "database-firewall-rules",
		Category:  "Databases NG",
		Usage:     "List firewall rules of a database",
		ArgsUsage: "database-id",
		Flags:     []cli.Flag{databaseFlag()},
		Description: CommandDescription{
			Description: "List all firewall rules of a database next generation",
			Examples: []string{
				"scalingo database-firewall-rules my-db-id",
				"scalingo --database my-db database-firewall-rules",
			},
			SeeAlso: []string{"database-firewall-rules-add", "database-firewall-rules-remove", "database-firewall-managed-ranges"},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			var databaseID, addonID string
			var err error

			databaseID = c.Args().First()
			if databaseID == "" {
				databaseID, addonID = detect.GetCurrentDatabase(ctx, c)
			} else {
				addonID, err = detect.GetAddonIDFromDatabase(ctx, databaseID)
				if err != nil {
					errorQuit(ctx, err)
				}
			}

			if databaseID == "" {
				io.Error("Please provide a database ID or use --database flag")
				return cli.ShowCommandHelp(ctx, c, "database-firewall-rules")
			}

			utils.CheckForConsent(ctx, databaseID, utils.ConsentTypeDBs)

			err = dbng.FirewallRulesList(ctx, databaseID, addonID)
			if err != nil {
				errorQuit(ctx, err)
			}

			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "database-firewall-rules")
			_ = autocomplete.DatabasesNgListAutoComplete(ctx)
		},
	}

	databaseFirewallRulesAddCommand = cli.Command{
		Name:      "database-firewall-rules-add",
		Category:  "Databases NG",
		Usage:     "Add a firewall rule to a database",
		ArgsUsage: "database-id",
		Flags: []cli.Flag{
			databaseFlag(),
			&cli.StringFlag{Name: "cidr", Usage: "CIDR range (e.g., 203.0.113.0/24)"},
			&cli.StringFlag{Name: "label", Usage: "Label for the custom rule"},
			&cli.StringFlag{Name: "managed-range", Usage: "Managed range ID"},
		},
		Description: CommandDescription{
			Description: "Add a firewall rule to a database next generation. Either --cidr or --managed-range must be specified, but not both.",
			Examples: []string{
				"scalingo database-firewall-rules-add my-db-id --cidr 203.0.113.0/24 --label \"Office network\"",
				"scalingo --database my-db database-firewall-rules-add --managed-range mr-scalingo-osc-fr1",
			},
			SeeAlso: []string{"database-firewall-rules", "database-firewall-rules-remove", "database-firewall-managed-ranges"},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			var databaseID, addonID string
			var err error

			databaseID = c.Args().First()
			if databaseID == "" {
				databaseID, addonID = detect.GetCurrentDatabase(ctx, c)
			} else {
				addonID, err = detect.GetAddonIDFromDatabase(ctx, databaseID)
				if err != nil {
					errorQuit(ctx, err)
				}
			}

			if databaseID == "" {
				io.Error("Please provide a database ID or use --database flag")
				return cli.ShowCommandHelp(ctx, c, "database-firewall-rules-add")
			}

			cidr := c.String("cidr")
			label := c.String("label")
			managedRange := c.String("managed-range")

			if cidr == "" && managedRange == "" {
				io.Error("Either --cidr or --managed-range must be specified")
				return cli.ShowCommandHelp(ctx, c, "database-firewall-rules-add")
			}

			if cidr != "" && managedRange != "" {
				io.Error("Cannot specify both --cidr and --managed-range")
				return cli.ShowCommandHelp(ctx, c, "database-firewall-rules-add")
			}

			utils.CheckForConsent(ctx, databaseID, utils.ConsentTypeDBs)

			var params scalingo.FirewallRuleCreateParams
			if cidr != "" {
				params = scalingo.FirewallRuleCreateParams{
					Type:  scalingo.FirewallRuleTypeCustomRange,
					CIDR:  cidr,
					Label: label,
				}
			} else {
				params = scalingo.FirewallRuleCreateParams{
					Type:    scalingo.FirewallRuleTypeManagedRange,
					RangeID: managedRange,
				}
			}

			err = dbng.FirewallRulesAdd(ctx, databaseID, addonID, params)
			if err != nil {
				errorQuit(ctx, err)
			}

			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "database-firewall-rules-add")
			_ = autocomplete.DatabasesNgListAutoComplete(ctx)
		},
	}

	databaseFirewallRulesRemoveCommand = cli.Command{
		Name:      "database-firewall-rules-remove",
		Category:  "Databases NG",
		Usage:     "Remove a firewall rule from a database",
		ArgsUsage: "database-id rule-id",
		Flags:     []cli.Flag{databaseFlag()},
		Description: CommandDescription{
			Description: "Remove a firewall rule from a database next generation",
			Examples: []string{
				"scalingo database-firewall-rules-remove my-db-id rule-id",
				"scalingo --database my-db database-firewall-rules-remove rule-id",
			},
			SeeAlso: []string{"database-firewall-rules", "database-firewall-rules-add", "database-firewall-managed-ranges"},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			var databaseID, addonID, ruleID string
			var err error

			args := c.Args().Slice()
			if len(args) == 2 {
				// Both database-id and rule-id provided as positional args
				databaseID = args[0]
				ruleID = args[1]
				addonID, err = detect.GetAddonIDFromDatabase(ctx, databaseID)
				if err != nil {
					errorQuit(ctx, err)
				}
			} else if len(args) == 1 {
				// Only rule-id provided, database from --database flag
				databaseID, addonID = detect.GetCurrentDatabase(ctx, c)
				ruleID = args[0]
			} else {
				io.Error("Please provide a rule ID")
				return cli.ShowCommandHelp(ctx, c, "database-firewall-rules-remove")
			}

			if databaseID == "" {
				io.Error("Please provide a database ID or use --database flag")
				return cli.ShowCommandHelp(ctx, c, "database-firewall-rules-remove")
			}

			utils.CheckForConsent(ctx, databaseID, utils.ConsentTypeDBs)

			err = dbng.FirewallRulesRemove(ctx, databaseID, addonID, ruleID)
			if err != nil {
				errorQuit(ctx, err)
			}

			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "database-firewall-rules-remove")
			_ = autocomplete.DatabasesNgListAutoComplete(ctx)
		},
	}

	databaseFirewallManagedRangesCommand = cli.Command{
		Name:      "database-firewall-managed-ranges",
		Category:  "Databases NG",
		Usage:     "List available managed ranges for a database",
		ArgsUsage: "database-id",
		Flags:     []cli.Flag{databaseFlag()},
		Description: CommandDescription{
			Description: "List all available managed ranges for firewall rules of a database next generation",
			Examples: []string{
				"scalingo database-firewall-managed-ranges my-db-id",
				"scalingo --database my-db database-firewall-managed-ranges",
			},
			SeeAlso: []string{"database-firewall-rules", "database-firewall-rules-add", "database-firewall-rules-remove"},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			var databaseID, addonID string
			var err error

			databaseID = c.Args().First()
			if databaseID == "" {
				databaseID, addonID = detect.GetCurrentDatabase(ctx, c)
			} else {
				addonID, err = detect.GetAddonIDFromDatabase(ctx, databaseID)
				if err != nil {
					errorQuit(ctx, err)
				}
			}

			if databaseID == "" {
				io.Error("Please provide a database ID or use --database flag")
				return cli.ShowCommandHelp(ctx, c, "database-firewall-managed-ranges")
			}

			utils.CheckForConsent(ctx, databaseID, utils.ConsentTypeDBs)

			err = dbng.FirewallManagedRangesList(ctx, databaseID, addonID)
			if err != nil {
				errorQuit(ctx, err)
			}

			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "database-firewall-managed-ranges")
			_ = autocomplete.DatabasesNgListAutoComplete(ctx)
		},
	}
)
