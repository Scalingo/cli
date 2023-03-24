package notifiers

import (
	"context"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
)

func Update(ctx context.Context, app, ID string, params ProvisionParams) error {
	if app == "" {
		return errgo.New("no app defined")
	}

	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	notifier, err := c.NotifierByID(ctx, app, ID)
	if err != nil {
		return errgo.Notef(err, "fail to get notifier from server")
	}

	eventTypes, err := c.EventTypesList(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to list event types")
	}

	// If there is no selected events, keep the existing ones
	if len(params.SelectedEventNames) == 0 {
		params.SelectedEventIDs = append(params.SelectedEventIDs, notifier.SelectedEventIDs...)
	} else {
		for _, name := range params.SelectedEventNames {
			for _, t := range eventTypes {
				if t.Name == name {
					params.SelectedEventIDs = append(params.SelectedEventIDs, t.ID)
					break
				}
			}
		}
	}

	if len(params.CollaboratorUsernames) > 0 {
		params.UserIDs, err = collaboratorUserIDs(ctx, c, app, params.CollaboratorUsernames)
		if err != nil {
			return errgo.Notef(err, "invalid collaborator usernames")
		}
	}

	baseNotifier, err := c.NotifierUpdate(ctx, app, ID, params.NotifierParams)
	if err != nil {
		return errgo.Notef(err, "fail to update notifier")
	}
	detailedNotifier := baseNotifier.Specialize()

	displayDetails(detailedNotifier, eventTypes)
	io.Info()
	io.Status("Notifier updated")
	return nil
}
