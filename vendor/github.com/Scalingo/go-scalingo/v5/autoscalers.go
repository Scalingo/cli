package scalingo

import (
	"context"

	"gopkg.in/errgo.v1"
)

type AutoscalersService interface {
	AutoscalersList(ctx context.Context, app string) ([]Autoscaler, error)
	AutoscalerAdd(ctx context.Context, app string, params AutoscalerAddParams) (*Autoscaler, error)
	AutoscalerRemove(ctx context.Context, app string, id string) error
}

var _ AutoscalersService = (*Client)(nil)

type Autoscaler struct {
	ID            string  `json:"id"`
	AppID         string  `json:"app_id"`
	ContainerType string  `json:"container_type"`
	Metric        string  `json:"metric"`
	Target        float64 `json:"target"`
	MinContainers int     `json:"min_containers"`
	MaxContainers int     `json:"max_containers"`
	Disabled      bool    `json:"disabled"`
}

type AutoscalersRes struct {
	Autoscalers []Autoscaler `json:"autoscalers"`
}

type AutoscalerRes struct {
	Autoscaler Autoscaler `json:"autoscaler"`
}

func (c *Client) AutoscalersList(ctx context.Context, app string) ([]Autoscaler, error) {
	var autoscalersRes AutoscalersRes
	err := c.ScalingoAPI().SubresourceList(ctx, "apps", app, "autoscalers", nil, &autoscalersRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return autoscalersRes.Autoscalers, nil
}

type AutoscalerAddParams struct {
	ContainerType string  `json:"container_type"`
	Metric        string  `json:"metric"`
	Target        float64 `json:"target"`
	MinContainers int     `json:"min_containers"`
	MaxContainers int     `json:"max_containers"`
}

func (c *Client) AutoscalerAdd(ctx context.Context, app string, params AutoscalerAddParams) (*Autoscaler, error) {
	var autoscalerRes AutoscalerRes
	err := c.ScalingoAPI().SubresourceAdd(ctx, "apps", app, "autoscalers", AutoscalerRes{
		Autoscaler: Autoscaler{
			ContainerType: params.ContainerType,
			Metric:        params.Metric,
			Target:        params.Target,
			MinContainers: params.MinContainers,
			MaxContainers: params.MaxContainers,
		},
	}, &autoscalerRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return &autoscalerRes.Autoscaler, nil
}

func (c *Client) AutoscalerShow(ctx context.Context, app, id string) (*Autoscaler, error) {
	var autoscalerRes AutoscalerRes
	err := c.ScalingoAPI().SubresourceGet(ctx, "apps", app, "autoscalers", id, nil, &autoscalerRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return &autoscalerRes.Autoscaler, nil
}

type AutoscalerUpdateParams struct {
	Metric        *string  `json:"metric,omitempty"`
	Target        *float64 `json:"target,omitempty"`
	MinContainers *int     `json:"min_containers,omitempty"`
	MaxContainers *int     `json:"max_containers,omitempty"`
	Disabled      *bool    `json:"disabled,omitempty"`
}

func (c *Client) AutoscalerUpdate(ctx context.Context, app, id string, params AutoscalerUpdateParams) (*Autoscaler, error) {
	var autoscalerRes AutoscalerRes
	err := c.ScalingoAPI().SubresourceUpdate(ctx, "apps", app, "autoscalers", id, params, &autoscalerRes)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return &autoscalerRes.Autoscaler, nil
}

func (c *Client) AutoscalerRemove(ctx context.Context, app, id string) error {
	err := c.ScalingoAPI().SubresourceDelete(ctx, "apps", app, "autoscalers", id)
	if err != nil {
		return errgo.Mask(err)
	}
	return nil
}
