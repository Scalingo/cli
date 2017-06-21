package notifiers

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/debug"
	"github.com/Scalingo/cli/io"
	scalingo "github.com/Scalingo/go-scalingo"
	"gopkg.in/errgo.v1"
)

func Provision(app, platformName string, params scalingo.NotifierCreateParams) error {
	if app == "" {
		return errgo.New("no app defined")
	}

	debug.Printf("[Provision] params: %+v", params)

	if len(params.SelectedEvents) >= 1 && params.SelectedEvents[0] == "" {
		params.SelectedEvents = nil
	}
	params.PlatformID = "593ac2d22664cd0001be2d0c"

	c := config.ScalingoClient()
	_, err := c.NotifierProvision(app, platformName, params)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	io.Status("Notifier have been created.")
	return nil
}
