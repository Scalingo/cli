package dbdr

import (
	"context"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-utils/errors/v3"
)

func DatabaseEndpointsList(ctx context.Context, databaseIdentifier string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	database, err := databaseByNameOrID(ctx, c, databaseIdentifier)
	if err != nil {
		return errors.Wrap(ctx, err, "show database")
	}

	endpoints, err := c.Preview().DatabaseEndpointsList(ctx, database.ID, scalingo.DatabaseEndpointsListParams{})
	if err != nil {
		return errors.Wrap(ctx, err, "list database endpoints")
	}

	if len(endpoints) == 0 {
		io.Status("No endpoint configured for this database.")
		return nil
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.Header([]string{"ID", "Type", "Hostname", "Port"})

	for _, endpoint := range endpoints {
		_ = t.Append([]string{
			endpoint.ID,
			string(endpoint.Type),
			endpoint.Hostname,
			strconv.Itoa(endpoint.Port),
		})
	}

	_ = t.Render()

	return nil
}

func DatabaseNetPeeringsList(ctx context.Context, databaseIdentifier string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	database, err := databaseByNameOrID(ctx, c, databaseIdentifier)
	if err != nil {
		return errors.Wrap(ctx, err, "show database")
	}

	netPeerings, err := c.Preview().DatabaseNetPeeringsList(ctx, database.ID)
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
			string(netPeering.Status),
			netPeering.OutscaleNetPeeringID,
			netPeering.OutscaleSourceNetID,
			netPeering.OutscaleSourceNetIPRange,
			netPeering.OutscaleSourceAccountID,
		})
	}

	_ = t.Render()

	return nil
}

func DatabaseNetPeeringsAdd(ctx context.Context, databaseIdentifier string, params scalingo.DatabaseNetPeeringCreateParams) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	database, err := databaseByNameOrID(ctx, c, databaseIdentifier)
	if err != nil {
		return errors.Wrap(ctx, err, "show database")
	}

	netPeering, err := c.Preview().DatabaseNetPeeringCreate(ctx, database.ID, params)
	if err != nil {
		return errors.Wrap(ctx, err, "add database net peering")
	}

	io.Statusf("Net peering '%s' creation has been initiated for database '%s'.\n", netPeering.ID, database.ID)
	io.Warning("Expect some delay for the net peering to be applied.")

	return nil
}

func DatabaseNetPeeringsRemove(ctx context.Context, databaseIdentifier, netPeeringID string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	database, err := databaseByNameOrID(ctx, c, databaseIdentifier)
	if err != nil {
		return errors.Wrap(ctx, err, "show database")
	}

	err = c.Preview().DatabaseNetPeeringDestroy(ctx, database.ID, netPeeringID)
	if err != nil {
		return errors.Wrap(ctx, err, "remove database net peering")
	}

	io.Statusf("Net peering '%s' has been removed from database '%s'.\n", netPeeringID, database.ID)
	io.Warning("Expect some delay for the net peering removal to be applied.")

	return nil
}

func DatabaseNetworkConfigurationShow(ctx context.Context, databaseIdentifier string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	database, err := databaseByNameOrID(ctx, c, databaseIdentifier)
	if err != nil {
		return errors.Wrap(ctx, err, "show database")
	}

	networkConfiguration, err := c.Preview().DatabaseNetworkConfigurationShow(ctx, database.ID)
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
