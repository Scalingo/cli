package notifications

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"gopkg.in/errgo.v1"
)

func Provision(app, webHookURL string) error {
	if app == "" {
		return errgo.New("no app defined")
	} else if webHookURL == "" {
		return errgo.New("no url defined")
	}

	c := config.ScalingoClient()
	params, err := c.NotificationProvision(app, webHookURL)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	io.Status("Notifications to", webHookURL, "have been created.")
	return nil
}
