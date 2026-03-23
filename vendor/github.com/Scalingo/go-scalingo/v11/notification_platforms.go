package scalingo

import (
	"context"

	"github.com/Scalingo/go-scalingo/v11/debug"
	"github.com/Scalingo/go-scalingo/v11/http"
	"github.com/Scalingo/go-utils/errors/v3"
)

type NotificationPlatformsService interface {
	NotificationPlatformsList(context.Context) ([]*NotificationPlatform, error)
	NotificationPlatformByName(ctx context.Context, name string) ([]*NotificationPlatform, error)
}

var _ NotificationPlatformsService = (*Client)(nil)

type NotificationPlatform struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	DisplayName       string   `json:"display_name"`
	LogoURL           string   `json:"logo_url"`
	Description       string   `json:"description"`
	AvailableEventIDs []string `json:"available_event_ids"`
}

type PlatformRes struct {
	NotificationPlatform *NotificationPlatform `json:"notification_platform"`
}

type PlatformsRes struct {
	NotificationPlatforms []*NotificationPlatform `json:"notification_platforms"`
}

func (c *Client) NotificationPlatformsList(ctx context.Context) ([]*NotificationPlatform, error) {
	var platformsRes PlatformsRes
	req := &http.APIRequest{
		NoAuth:   true,
		Endpoint: "/notification_platforms",
	}
	err := c.ScalingoAPI().DoRequest(ctx, req, &platformsRes)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "list notification platforms")
	}

	return platformsRes.NotificationPlatforms, nil
}

func (c *Client) NotificationPlatformByName(ctx context.Context, name string) ([]*NotificationPlatform, error) {
	debug.Printf("[NotificationPlatformByName] name: %s", name)
	var platformsRes PlatformsRes
	req := &http.APIRequest{
		NoAuth:   true,
		Endpoint: "/notification_platforms/search",
		Params:   map[string]string{"name": name},
	}
	err := c.ScalingoAPI().DoRequest(ctx, req, &platformsRes)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "search notification platforms by name")
	}

	return platformsRes.NotificationPlatforms, nil
}
