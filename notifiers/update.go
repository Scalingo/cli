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

	c := config.ScalingoClient()
	notifier, err := c.NotifierByID(app, ID)
	if err != nil {
		return errgo.Notef(err, "fail to get notifier from server")
	}

	// If there is no selected events, keep the existing ones
	if len(params.SelectedEvents) == 1 && params.SelectedEvents[0] == "" {
		params.SelectedEvents = []string{}
		for _, e := range notifier.SelectedEvents {
			params.SelectedEvents = append(params.SelectedEvents, e.Name)
		}
	}

	baseNotifier, err := c.NotifierUpdate(app, ID, string(notifier.GetType()), params)
	if err != nil {
		return errgo.Notef(err, "fail to update notifier")
	}
	detailedNotifier := baseNotifier.Specialize()

	displayDetails(detailedNotifier)
	io.Info()
	io.Status("Notifier updated")
	return nil
}
