package scalingo

import (
	"encoding/json"

	errgo "gopkg.in/errgo.v1"
)

type NotificationPlatform struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	DisplayName     string   `json:"display_name"`
	EventsAvailable []string `json:"events_available"`
}

type PlatformRes struct {
	NotificationPlatform *NotificationPlatform `json:"notification_platform"`
}

type PlatformsRes struct {
	NotificationPlatforms []*NotificationPlatform `json:"notification_platform"`
}

func (c *Client) NotificationPlatformsList() ([]*NotificationPlatform, error) {
	req := &APIRequest{
		Client:   c,
		NoAuth:   true,
		Endpoint: "/notification_platforms",
	}
	res, err := req.Do()
	if err != nil {
		return nil, errgo.Mask(err)
	}
	defer res.Body.Close()

	var params PlatformsRes
	err = json.NewDecoder(res.Body).Decode(&params)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	return params.NotificationPlatforms, nil
}

func (c *Client) NotificationPlatformByName(name string) (*NotificationPlatform, error) {
	req := &APIRequest{
		Client:   c,
		NoAuth:   true,
		Endpoint: "/notification_platforms/search",
		Params:   map[string]string{"name": name},
	}
	res, err := req.Do()
	if err != nil {
		return nil, errgo.Mask(err)
	}
	defer res.Body.Close()

	var params PlatformRes
	err = json.NewDecoder(res.Body).Decode(&params)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	return params.NotificationPlatform, nil
}
