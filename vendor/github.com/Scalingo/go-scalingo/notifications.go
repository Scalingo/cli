package scalingo

import "gopkg.in/errgo.v1"

type Notification struct {
	ID              string `json:"id"`
	Type            string `json:"type"`
	WebHookURL      string `json:"webhook_url"`
	Active          bool   `json:"active"`
}

type NotificationRes struct {
	Notification     Notification `json:"notification"`
	Message          string       `json:"message,omitempty"`
	Variables        []string     `json:"variables,omitempty"`
}

type NotificationsRes struct {
	Notifications []*Notification `json:"notifications"`
}

func (c *Client) NotificationsList(app string) ([]*Notification, error) {
	var notificationsRes NotificationsRes
	err := c.subresourceList(app, "notifications", nil, &notificationsRes)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	return notificationsRes.Notifications, nil
}

func (c *Client) NotificationProvision(app, webHookURL string) (NotificationRes, error) {
	var notificationRes NotificationRes
	err := c.subresourceAdd(app, "notifications", NotificationRes{Notification: Notification{WebHookURL: webHookURL}}, &notificationRes)
	if err != nil {
		return NotificationRes{}, errgo.Mask(err, errgo.Any)
	}
	return notificationRes, nil
}

func (c *Client) NotificationUpdate(app, ID, webHookURL string) (NotificationRes, error) {
	var notificationRes NotificationRes
	err := c.subresourceUpdate(app, "notifications", ID, NotificationRes{Notification: Notification{WebHookURL: webHookURL, Active: true}}, &notificationRes)
	if err != nil {
		return NotificationRes{}, errgo.Mask(err, errgo.Any)
	}
	return notificationRes, nil
}

func (c *Client) NotificationDestroy(app, ID string) error {
	return c.subresourceDelete(app, "notifications", ID)
}
