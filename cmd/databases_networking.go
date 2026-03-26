package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/dbng"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v11"
)

var (
	databaseNetPeeringsListCommand = cli.Command{
		Name:      "database-net-peerings",
		Category:  "Databases DR",
		Usage:     "List net peerings of a database",
		ArgsUsage: "database-id",
		Flags:     []cli.Flag{databaseFlag()},
		Description: CommandDescription{
			Description: "List all net peerings of a database",
			Examples: []string{
				"scalingo database-net-peerings my-db-id",
				"scalingo --database my-db database-net-peerings",
			},
			SeeAlso: []string{"database-net-peerings-add", "database-net-peerings-remove", "database-network-configuration"},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			databaseID := detect.ExtractDatabaseNameFromCommandLineOrEnv(c)
			if databaseID == "" {
				return cli.ShowCommandHelp(ctx, c, "database-net-peerings")
			}
			utils.CheckForConsent(ctx, databaseID, utils.ConsentTypeDBs)

			err := dbng.DatabaseNetPeeringsList(ctx, databaseID)
			if err != nil {
				errorQuit(ctx, err)
			}

			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "database-net-peerings")
			_ = autocomplete.DatabasesNgListAutoComplete(ctx)
		},
	}

	databaseNetPeeringsAddCommand = cli.Command{
		Name:      "database-net-peerings-add",
		Category:  "Databases DR",
		Usage:     "Add a net peering to a database",
		ArgsUsage: "database-id",
		Flags: []cli.Flag{
			databaseFlag(),
			&cli.StringFlag{
				Name:     "outscale-net-peering-id",
				Usage:    "Outscale net peering ID",
				Required: true,
			},
		},
		Description: CommandDescription{
			Description: "Initiate the creation of a net peering for a database",
			Examples: []string{
				"scalingo database-net-peerings-add my-db-id --outscale-net-peering-id pcx-123456789",
				"scalingo --database my-db database-net-peerings-add --outscale-net-peering-id pcx-123456789",
			},
			SeeAlso: []string{"database-net-peerings", "database-net-peerings-remove", "database-network-configuration"},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			databaseID := detect.ExtractDatabaseNameFromCommandLineOrEnv(c)
			if databaseID == "" {
				return cli.ShowCommandHelp(ctx, c, "database-net-peerings-add")
			}

			utils.CheckForConsent(ctx, databaseID, utils.ConsentTypeDBs)

			params := scalingo.DatabaseNetPeeringCreateParams{
				OutscaleNetPeeringID: c.String("outscale-net-peering-id"),
			}

			err := dbng.DatabaseNetPeeringsAdd(ctx, databaseID, params)
			if err != nil {
				errorQuit(ctx, err)
			}

			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "database-net-peerings-add")
			_ = autocomplete.DatabasesNgListAutoComplete(ctx)
		},
	}

	databaseNetPeeringsRemoveCommand = cli.Command{
		Name:      "database-net-peerings-remove",
		Category:  "Databases DR",
		Usage:     "Remove a net peering from a database",
		ArgsUsage: "database-id net-peering-id",
		Flags: []cli.Flag{
			databaseFlag(),
			&cli.StringFlag{
				Name:    "net-peering",
				Aliases: []string{"np"},
				Usage:   "Net peering ID to remove",
			},
		},
		Description: CommandDescription{
			Description: "Delete an existing net peering from a database",
			Examples: []string{
				"scalingo database-net-peerings-remove my-db-id np-id",
				"scalingo --database my-db database-net-peerings-remove np-id",
			},
			SeeAlso: []string{"database-net-peerings", "database-net-peerings-add", "database-network-configuration"},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			databaseID := detect.ExtractDatabaseNameFromCommandLineOrEnv(c)
			if databaseID == "" {
				return cli.ShowCommandHelp(ctx, c, "database-net-peerings-remove")
			}

			utils.CheckForConsent(ctx, databaseID, utils.ConsentTypeDBs)

			err := dbng.DatabaseNetPeeringsRemove(ctx, databaseID, c.String("net-peering"))
			if err != nil {
				errorQuit(ctx, err)
			}

			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "database-net-peerings-remove")
			_ = autocomplete.DatabasesNgListAutoComplete(ctx)
		},
	}

	databaseNetworkConfigurationShowCommand = cli.Command{
		Name:      "database-network-configuration",
		Category:  "Databases DR",
		Usage:     "Show network configuration of a database",
		ArgsUsage: "database-id",
		Flags:     []cli.Flag{databaseFlag()},
		Description: CommandDescription{
			Description: "Get network configuration of a database",
			Examples: []string{
				"scalingo database-network-configuration my-db-id",
				"scalingo --database my-db database-network-configuration",
			},
			SeeAlso: []string{"database-net-peerings", "database-net-peerings-add", "database-net-peerings-remove"},
		}.Render(),
		Action: func(ctx context.Context, c *cli.Command) error {
			databaseID := detect.ExtractDatabaseNameFromCommandLineOrEnv(c)
			if databaseID == "" {
				return cli.ShowCommandHelp(ctx, c, "database-network-configuration")
			}

			utils.CheckForConsent(ctx, databaseID, utils.ConsentTypeDBs)

			err := dbng.DatabaseNetworkConfigurationShow(ctx, databaseID)
			if err != nil {
				errorQuit(ctx, err)
			}

			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "database-network-configuration")
			_ = autocomplete.DatabasesNgListAutoComplete(ctx)
		},
	}
)
