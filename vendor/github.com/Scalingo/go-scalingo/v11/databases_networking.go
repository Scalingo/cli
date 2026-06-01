package scalingo

import (
	"context"
	"time"

	"github.com/Scalingo/go-utils/errors/v3"
)

type DatabaseEndpointType string

const (
	DatabaseEndpointTypePublicRW         DatabaseEndpointType = "public-rw"
	DatabaseEndpointTypePrivatePeeringRW DatabaseEndpointType = "private-peering-rw"
)

type DatabaseEndpoint struct {
	ID         string               `json:"id"`
	DatabaseID string               `json:"database_id"`
	Hostname   string               `json:"hostname"`
	Port       int                  `json:"port"`
	Type       DatabaseEndpointType `json:"type"`
}

type DatabaseEndpointsResponse struct {
	Endpoints []DatabaseEndpoint `json:"endpoints"`
}

type DatabaseNetPeeringStatus string

const (
	DatabaseNetPeeringStatusActive  DatabaseNetPeeringStatus = "active"
	DatabaseNetPeeringStatusDeleted DatabaseNetPeeringStatus = "deleted"
)

type DatabaseNetPeering struct {
	ID                       string                   `json:"id"`
	DatabaseID               string                   `json:"database_id"`
	Status                   DatabaseNetPeeringStatus `json:"status"`
	OutscaleNetPeeringID     string                   `json:"outscale_net_peering_id"`
	OutscaleSourceNetID      string                   `json:"outscale_source_net_id"`
	OutscaleSourceNetIPRange string                   `json:"outscale_source_net_ip_range"`
	OutscaleSourceAccountID  string                   `json:"outscale_source_account_id"`
	CreatedAt                time.Time                `json:"created_at"`
	UpdatedAt                time.Time                `json:"updated_at"`
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

func (c *PreviewClient) DatabaseEndpointsList(ctx context.Context, databaseID string) ([]DatabaseEndpoint, error) {
	var res DatabaseEndpointsResponse

	err := c.parent.ScalingoAPI().SubresourceList(ctx, databasesResource, databaseID, databaseEndpointsResource, nil, &res)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "list database endpoints")
	}

	return res.Endpoints, nil
}

func (c *PreviewClient) DatabaseNetPeeringCreate(ctx context.Context, databaseID string, params DatabaseNetPeeringCreateParams) (DatabaseNetPeering, error) {
	var res DatabaseNetPeeringResponse

	err := c.parent.ScalingoAPI().SubresourceAdd(ctx, databasesResource, databaseID, netPeeringsResource, params, &res)
	if err != nil {
		return DatabaseNetPeering{}, errors.Wrap(ctx, err, "create database net peering")
	}

	return res.NetPeering, nil
}

func (c *PreviewClient) DatabaseNetPeeringsList(ctx context.Context, databaseID string) ([]DatabaseNetPeering, error) {
	var res DatabaseNetPeeringsResponse

	err := c.parent.ScalingoAPI().SubresourceList(ctx, databasesResource, databaseID, netPeeringsResource, nil, &res)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "list database net peerings")
	}

	return res.NetPeerings, nil
}

func (c *PreviewClient) DatabaseNetPeeringShow(ctx context.Context, databaseID, netPeeringID string) (DatabaseNetPeering, error) {
	var res DatabaseNetPeeringResponse

	err := c.parent.ScalingoAPI().SubresourceGet(ctx, databasesResource, databaseID, netPeeringsResource, netPeeringID, nil, &res)
	if err != nil {
		return DatabaseNetPeering{}, errors.Wrap(ctx, err, "show database net peering")
	}

	return res.NetPeering, nil
}

func (c *PreviewClient) DatabaseNetPeeringDestroy(ctx context.Context, databaseID, netPeeringID string) error {
	err := c.parent.ScalingoAPI().SubresourceDelete(ctx, databasesResource, databaseID, netPeeringsResource, netPeeringID)
	if err != nil {
		return errors.Wrap(ctx, err, "destroy database net peering")
	}

	return nil
}

func (c *PreviewClient) DatabaseNetworkConfigurationShow(ctx context.Context, databaseID string) (DatabaseNetworkConfiguration, error) {
	var res DatabaseNetworkConfigurationResponse

	err := c.parent.ScalingoAPI().SubresourceGetSingleton(ctx, databasesResource, databaseID, networkConfigurationResource, nil, &res)
	if err != nil {
		return DatabaseNetworkConfiguration{}, errors.Wrap(ctx, err, "show database network configuration")
	}

	return res.NetworkConfiguration, nil
}
