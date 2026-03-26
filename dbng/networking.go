package dbng

import (
	"context"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-utils/errors/v3"
)

func DatabaseNetPeeringsList(ctx context.Context, databaseID string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	netPeerings, err := c.Preview().DatabaseNetPeeringsList(ctx, databaseID)
	if err != nil {
		return errors.Wrap(ctx, err, "list database net peerings")
	}

	if len(netPeerings) == 0 {
		io.Status("No net peering configured for this database.")
		return nil
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.Header([]string{"ID", "Status", "Outscale Net Peering ID", "Source Net ID", "Source Net IP Range", "Source Account ID"})

	for _, netPeering := range netPeerings {
		_ = t.Append([]string{
			netPeering.ID,
			netPeering.Status,
			netPeering.OutscaleNetPeeringID,
			netPeering.OutscaleSourceNetID,
			netPeering.OutscaleSourceNetIPRange,
			netPeering.OutscaleSourceAccountID,
		})
	}

	_ = t.Render()

	return nil
}

func DatabaseNetPeeringsAdd(ctx context.Context, databaseID string, params scalingo.DatabaseNetPeeringCreateParams) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	netPeering, err := c.Preview().DatabaseNetPeeringCreate(ctx, databaseID, params)
	if err != nil {
		return errors.Wrap(ctx, err, "add database net peering")
	}

	io.Statusf("Net peering '%s' creation has been initiated for database '%s'.\n", netPeering.ID, databaseID)
	io.Warning("Expect some delay for the net peering to be applied.")

	return nil
}

func DatabaseNetPeeringsRemove(ctx context.Context, databaseID, netPeeringID string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	err = c.Preview().DatabaseNetPeeringDestroy(ctx, databaseID, netPeeringID)
	if err != nil {
		return errors.Wrap(ctx, err, "remove database net peering")
	}

	io.Statusf("Net peering '%s' has been removed from database '%s'.\n", netPeeringID, databaseID)
	io.Warning("Expect some delay for the net peering removal to be applied.")

	return nil
}

func DatabaseNetworkConfigurationShow(ctx context.Context, databaseID string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	networkConfiguration, err := c.Preview().DatabaseNetworkConfigurationShow(ctx, databaseID)
	if err != nil {
		return errors.Wrap(ctx, err, "show database network configuration")
	}

	t := tablewriter.NewWriter(os.Stdout)
	_ = t.Bulk([][]string{
		{"Outscale account ID", networkConfiguration.OutscaleAccountID},
		{"Outscale net ID", networkConfiguration.OutscaleNetID},
		{"IP range", networkConfiguration.IPRange},
	})
	_ = t.Render()

	return nil
}
