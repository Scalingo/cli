package notifications

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"gopkg.in/errgo.v1"
)

func Update(app, ID, webHookURL string) error {
	if app == "" {
		return errgo.New("no app defined")
	} else if webHookURL == "" {
		return errgo.New("no url defined")
	}

	_, err := checkNotificationExist(app, ID)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	c := config.ScalingoClient()
	_, err = c.NotificationUpdate(app, ID, webHookURL)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	io.Status("Notifications are now sent to", webHookURL)
	return nil
}
