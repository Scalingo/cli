package scalingo

import (
	"context"
	"time"

	httpclient "github.com/Scalingo/go-scalingo/v11/http"
	"github.com/Scalingo/go-utils/errors/v3"
)

type DatabaseNetPeering struct {
	ID                       string    `json:"id"`
	DatabaseID               string    `json:"database_id"`
	Status                   string    `json:"status"`
	OutscaleNetPeeringID     string    `json:"outscale_net_peering_id"`
	OutscaleSourceNetID      string    `json:"outscale_source_net_id"`
	OutscaleSourceNetIPRange string    `json:"outscale_source_net_ip_range"`
	OutscaleSourceAccountID  string    `json:"outscale_source_account_id"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
}

type DatabaseNetPeeringCreateParams struct {
	OutscaleNetPeeringID string `json:"outscale_net_peering_id"`
}

type DatabaseNetPeeringResponse struct {
	NetPeering DatabaseNetPeering `json:"net_peering"`
}

type DatabaseNetPeeringsResponse struct {
	NetPeerings []DatabaseNetPeering `json:"net_peerings"`
}

type DatabaseNetworkConfiguration struct {
	OutscaleAccountID string `json:"outscale_account_id"`
	OutscaleNetID     string `json:"outscale_net_id"`
	IPRange           string `json:"ip_range"`
}

type DatabaseNetworkConfigurationResponse struct {
	NetworkConfiguration DatabaseNetworkConfiguration `json:"network_configuration"`
}

func (c *PreviewClient) DatabaseNetPeeringCreate(ctx context.Context, databaseID string, params DatabaseNetPeeringCreateParams) (DatabaseNetPeering, error) {
	var res DatabaseNetPeeringResponse

	err := c.parent.ScalingoAPI().SubresourceAdd(ctx, "databases", databaseID, "net_peerings", params, &res)
	if err != nil {
		return DatabaseNetPeering{}, errors.Wrap(ctx, err, "create database net peering")
	}

	return res.NetPeering, nil
}

func (c *PreviewClient) DatabaseNetPeeringsList(ctx context.Context, databaseID string) ([]DatabaseNetPeering, error) {
	var res DatabaseNetPeeringsResponse

	err := c.parent.ScalingoAPI().SubresourceList(ctx, "databases", databaseID, "net_peerings", nil, &res)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "list database net peerings")
	}

	return res.NetPeerings, nil
}

func (c *PreviewClient) DatabaseNetPeeringShow(ctx context.Context, databaseID, netPeeringID string) (DatabaseNetPeering, error) {
	var res DatabaseNetPeeringResponse

	err := c.parent.ScalingoAPI().SubresourceGet(ctx, "databases", databaseID, "net_peerings", netPeeringID, nil, &res)
	if err != nil {
		return DatabaseNetPeering{}, errors.Wrap(ctx, err, "show database net peering")
	}

	return res.NetPeering, nil
}

func (c *PreviewClient) DatabaseNetPeeringDestroy(ctx context.Context, databaseID, netPeeringID string) error {
	err := c.parent.ScalingoAPI().SubresourceDelete(ctx, "databases", databaseID, "net_peerings", netPeeringID)
	if err != nil {
		return errors.Wrap(ctx, err, "destroy database net peering")
	}

	return nil
}

func (c *PreviewClient) DatabaseNetworkConfigurationShow(ctx context.Context, databaseID string) (DatabaseNetworkConfiguration, error) {
	var res DatabaseNetworkConfigurationResponse

	req := &httpclient.APIRequest{
		Method:   "GET",
		Endpoint: "/databases/" + databaseID + "/network_configuration",
	}

	err := c.parent.ScalingoAPI().DoRequest(ctx, req, &res)
	if err != nil {
		return DatabaseNetworkConfiguration{}, errors.Wrap(ctx, err, "show database network configuration")
	}

	return res.NetworkConfiguration, nil
}
