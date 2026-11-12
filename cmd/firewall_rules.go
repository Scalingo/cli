package cmd

import (
	"context"
	"errors"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/dbdr"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v11"
)

var (
	databaseFirewallRulesListCommand = cli.Command{
		Name:      "database-firewall-rules",
		Category:  "Databases DR",
		Usage:     "List firewall rules of a database",
		ArgsUsage: "database-id-or-name",
		Flags:     []cli.Flag{databaseFlag()},
		Description: CommandDescription{
			Description: "List all firewall rules of a database",
			Examples: []string{
				"scalingo database-firewall-rules database_id_or_name",
				"scalingo --database my-db database-firewall-rules",
			},
			SeeAlso: []string{"database-firewall-rules-add", "database-firewall-rules-remove", "database-firewall-managed-ranges"},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			databaseIdentifier, err := databaseIdentifierFromArgs(c)
			if errors.Is(err, detect.ErrTooManyArguments) {
				io.Error(err)
				return cli.ShowCommandHelp(ctx, c, "database-firewall-rules")
			}
			if databaseIdentifier == "" {
				io.Error("Please provide a database ID or name, or use --database flag")
				return cli.ShowCommandHelp(ctx, c, "database-firewall-rules")
			}

			utils.CheckForConsent(ctx, databaseIdentifier, utils.ConsentTypeDBs)

			err = dbdr.FirewallRulesList(ctx, databaseIdentifier)
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
		Category:  "Databases DR",
		Usage:     "Add a firewall rule to a database",
		ArgsUsage: "database-id-or-name",
		Flags: []cli.Flag{
			databaseFlag(),
			&cli.StringFlag{Name: "cidr", Usage: "CIDR range (e.g., 203.0.113.0/24)"},
			&cli.StringFlag{Name: "label", Usage: "Label for the custom rule"},
			&cli.StringFlag{Name: "managed-range", Usage: "Managed range ID"},
		},
		Description: CommandDescription{
			Description: "Add a firewall rule to a database. Either --cidr or --managed-range must be specified, but not both.",
			Examples: []string{
				"scalingo database-firewall-rules-add database_id_or_name --cidr 203.0.113.0/24 --label \"Office network\"",
				"scalingo --database my-db database-firewall-rules-add --managed-range mr-scalingo-osc-fr1",
			},
			SeeAlso: []string{"database-firewall-rules", "database-firewall-rules-remove", "database-firewall-managed-ranges"},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			databaseIdentifier, err := databaseIdentifierFromArgs(c)
			if errors.Is(err, detect.ErrTooManyArguments) {
				io.Error(err)
				return cli.ShowCommandHelp(ctx, c, "database-firewall-rules-add")
			}
			if databaseIdentifier == "" {
				io.Error("Please provide a database ID or name, or use --database flag")
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

			if managedRange != "" && label != "" {
				io.Error("Cannot specify --label with --managed-range")
				return cli.ShowCommandHelp(ctx, c, "database-firewall-rules-add")
			}

			utils.CheckForConsent(ctx, databaseIdentifier, utils.ConsentTypeDBs)

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

			err = dbdr.FirewallRulesAdd(ctx, databaseIdentifier, params)
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
		Category:  "Databases DR",
		Usage:     "Remove a firewall rule from a database",
		ArgsUsage: "database-id-or-name rule-id",
		Flags:     []cli.Flag{databaseFlag()},
		Description: CommandDescription{
			Description: "Remove a firewall rule from a database",
			Examples: []string{
				"scalingo database-firewall-rules-remove database_id_or_name rule_id",
				"scalingo --database my-db database-firewall-rules-remove rule-id",
			},
			SeeAlso: []string{"database-firewall-rules", "database-firewall-rules-add", "database-firewall-managed-ranges"},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			var databaseIdentifier, ruleID string

			args := c.Args().Slice()
			switch len(args) {
			case 0:
				// No arguments provided
				io.Error("Missing required argument: rule-id")
				return cli.ShowCommandHelp(ctx, c, "database-firewall-rules-remove")
			case 1:
				// Only rule-id provided, database from --database flag
				databaseIdentifier = detect.ExtractDatabaseNameFromCommandLineOrEnv(c)
				ruleID = args[0]
			case 2:
				// Both database identifier and rule ID provided as positional args
				databaseIdentifier = args[0]
				ruleID = args[1]
			default:
				io.Error("Invalid number of arguments")
				return cli.ShowCommandHelp(ctx, c, "database-firewall-rules-remove")
			}

			if databaseIdentifier == "" {
				io.Error("Please provide a database ID or name, or use --database flag")
				return cli.ShowCommandHelp(ctx, c, "database-firewall-rules-remove")
			}

			utils.CheckForConsent(ctx, databaseIdentifier, utils.ConsentTypeDBs)

			err := dbdr.FirewallRulesRemove(ctx, databaseIdentifier, ruleID)
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
		Category:  "Databases DR",
		Usage:     "List available managed ranges for a database",
		ArgsUsage: "database-id-or-name",
		Flags:     []cli.Flag{databaseFlag()},
		Description: CommandDescription{
			Description: "List all available managed ranges for firewall rules of a database",
			Examples: []string{
				"scalingo database-firewall-managed-ranges database_id_or_name",
				"scalingo --database my-db database-firewall-managed-ranges",
			},
			SeeAlso: []string{"database-firewall-rules", "database-firewall-rules-add", "database-firewall-rules-remove"},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			databaseIdentifier, err := databaseIdentifierFromArgs(c)
			if errors.Is(err, detect.ErrTooManyArguments) {
				io.Error(err)
				return cli.ShowCommandHelp(ctx, c, "database-firewall-managed-ranges")
			}
			if databaseIdentifier == "" {
				io.Error("Please provide a database ID or name, or use --database flag")
				return cli.ShowCommandHelp(ctx, c, "database-firewall-managed-ranges")
			}

			utils.CheckForConsent(ctx, databaseIdentifier, utils.ConsentTypeDBs)

			err = dbdr.FirewallManagedRangesList(ctx, databaseIdentifier)
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

func databaseIdentifierFromArgs(c *cli.Command) (string, error) {
	switch c.Args().Len() {
	case 0:
		return detect.ExtractDatabaseNameFromCommandLineOrEnv(c), nil
	case 1:
		return c.Args().First(), nil
	default:
		return "", detect.ErrTooManyArguments
	}
}
