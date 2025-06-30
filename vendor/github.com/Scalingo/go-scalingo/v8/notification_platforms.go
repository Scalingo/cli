package scalingo

import (
	"context"
	"encoding/json"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/go-scalingo/v8/debug"
	"github.com/Scalingo/go-scalingo/v8/http"
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
	req := &http.APIRequest{
		NoAuth:   true,
		Endpoint: "/notification_platforms",
	}
	res, err := c.ScalingoAPI().Do(ctx, req)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	defer res.Body.Close()

	var response PlatformsRes
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	return response.NotificationPlatforms, nil
}

func (c *Client) NotificationPlatformByName(ctx context.Context, name string) ([]*NotificationPlatform, error) {
	debug.Printf("[NotificationPlatformByName] name: %s", name)
	req := &http.APIRequest{
		NoAuth:   true,
		Endpoint: "/notification_platforms/search",
		Params:   map[string]string{"name": name},
	}
	res, err := c.ScalingoAPI().Do(ctx, req)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	defer res.Body.Close()

	debug.Printf("[NotificationPlatformByName] response: %+v", res.Body)
	var response PlatformsRes
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	return response.NotificationPlatforms, nil
}
