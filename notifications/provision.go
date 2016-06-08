package notifications

import (
	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
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

	io.Status("Notifications to", webHookURL, "have been provisionned")
	io.Info("ID:", params.Notification.ID)
	if len(params.Variables) > 0 {
		io.Info("Modified variables:", params.Variables)
	}
	if len(params.Message) > 0 {
		io.Info("Message from notification provider:", params.Message)
	}
	return nil
}
