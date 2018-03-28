package scalingo

import "gopkg.in/errgo.v1"

type NotificationsService interface {
	NotificationsList(app string) ([]*Notification, error)
	NotificationProvision(app, webHookURL string) (NotificationRes, error)
	NotificationUpdate(app, ID, webHookURL string) (NotificationRes, error)
	NotificationDestroy(app, ID string) error
}

type NotificationsClient struct {
	subresourceClient
}

type Notification struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	WebHookURL string `json:"webhook_url"`
	Active     bool   `json:"active"`
}

type NotificationRes struct {
	Notification Notification `json:"notification"`
}

type NotificationsRes struct {
	Notifications []*Notification `json:"notifications"`
}

func (c *NotificationsClient) NotificationsList(app string) ([]*Notification, error) {
	var notificationsRes NotificationsRes
	err := c.subresourceList(app, "notifications", nil, &notificationsRes)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	return notificationsRes.Notifications, nil
}

func (c *NotificationsClient) NotificationProvision(app, webHookURL string) (NotificationRes, error) {
	var notificationRes NotificationRes
	err := c.subresourceAdd(app, "notifications", NotificationRes{Notification: Notification{WebHookURL: webHookURL}}, &notificationRes)
	if err != nil {
		return NotificationRes{}, errgo.Mask(err, errgo.Any)
	}
	return notificationRes, nil
}

func (c *NotificationsClient) NotificationUpdate(app, ID, webHookURL string) (NotificationRes, error) {
	var notificationRes NotificationRes
	err := c.subresourceUpdate(app, "notifications", ID, NotificationRes{Notification: Notification{WebHookURL: webHookURL, Active: true}}, &notificationRes)
	if err != nil {
		return NotificationRes{}, errgo.Mask(err, errgo.Any)
	}
	return notificationRes, nil
}

func (c *NotificationsClient) NotificationDestroy(app, ID string) error {
	return c.subresourceDelete(app, "notifications", ID)
}
