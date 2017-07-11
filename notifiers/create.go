package notifiers

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/debug"
	"github.com/Scalingo/cli/io"
	scalingo "github.com/Scalingo/go-scalingo"
	"gopkg.in/errgo.v1"
)

func Provision(app, platformName string, params scalingo.NotifierParams) error {
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
	platforms, err := c.NotificationPlatformByName(platformName)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	if len(platforms) <= 0 {
		return errgo.Newf("notification platform \"%s\" has not been found", platformName)
	}
	params.PlatformID = platforms[0].ID

	baseNotifier, err := c.NotifierProvision(app, platforms[0].Name, params)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	notifier := baseNotifier.Specialize()

	displayDetails(notifier)

	io.Info()
	io.Status("Notifier have been created.")
	return nil
}
