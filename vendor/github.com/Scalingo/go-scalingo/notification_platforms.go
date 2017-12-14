package scalingo

import (
	"encoding/json"

	"github.com/Scalingo/go-scalingo/debug"

	errgo "gopkg.in/errgo.v1"
)

type NotificationPlatform struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	DisplayName     string            `json:"display_name"`
	AvailableEvents []EventTypeStruct `json:"available_events"`
}

type PlatformRes struct {
	NotificationPlatform *NotificationPlatform `json:"notification_platform"`
}

type PlatformsRes struct {
	NotificationPlatforms []*NotificationPlatform `json:"notification_platforms"`
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

	var response PlatformsRes
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	return response.NotificationPlatforms, nil
}

func (c *Client) NotificationPlatformByName(name string) ([]*NotificationPlatform, error) {
	debug.Printf("[NotificationPlatformByName] name: %s", name)
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

	debug.Printf("[NotificationPlatformByName] reponse: %+v", res.Body)
	var response PlatformsRes
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}

	return response.NotificationPlatforms, nil
}
