package notifiers

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	scalingo "github.com/Scalingo/go-scalingo"
	"gopkg.in/errgo.v1"
)

func Update(app, ID string, params scalingo.NotifierParams) error {
	if app == "" {
		return errgo.New("no app defined")
	}
	if len(params.SelectedEvents) >= 1 && params.SelectedEvents[0] == "" {
		params.SelectedEvents = nil
	}

	c := config.ScalingoClient()
	notifier, err := c.NotifierByID(app, ID)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	_, err = c.NotifierUpdate(app, ID, string(notifier.GetType()), params)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	io.Status("Notifier updated")
	return nil
}
