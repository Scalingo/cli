package notifiers

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/debug"
	"github.com/Scalingo/cli/io"
	scalingo "github.com/Scalingo/go-scalingo"
	"gopkg.in/errgo.v1"
)

func Provision(app, platformName string, params scalingo.NotifierCreateParams) error {
	debug.Printf("[Provision] params: %+v", params)

	if app == "" {
		return errgo.New("no app defined")
	}
	if platformName == "" {
		return errgo.New("no platform defined")
	}
	if len(params.SelectedEvents) >= 1 && params.SelectedEvents[0] == "" {
		params.SelectedEvents = nil
	}

	c := config.ScalingoClient()
	platform, err := c.NotificationPlatformByName(platformName)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	params.PlatformID = platform.ID

	_, err = c.NotifierProvision(app, platform.Name, params)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	io.Status("Notifier have been created.")
	return nil
}
