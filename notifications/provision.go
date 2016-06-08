package notifications

import (
	"errors"

	"gopkg.in/errgo.v1"
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
		return errors.New("Notification " + webHookURL + " can't be enabled because you have already an other notification of the same type enabled!")
	}

	io.Status("Notifications to", webHookURL, "have been provisionned")
	if len(params.Variables) > 0 {
		io.Info("Modified variables:", params.Variables)
	}
	if len(params.Message) > 0 {
		io.Info("Message from notification provider:", params.Message)
	}
	return nil
}
