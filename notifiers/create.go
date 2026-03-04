package notifiers

import (
	"context"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v10"
	"github.com/Scalingo/go-scalingo/v10/debug"
	"github.com/Scalingo/go-utils/errors/v3"
)

type ProvisionParams struct {
	CollaboratorUsernames []string
	SelectedEventNames    []string
	scalingo.NotifierParams
}

func Provision(ctx context.Context, app, platformName string, params ProvisionParams) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}
	debug.Printf("[Provision] params: %+v", params)

	if app == "" {
		return errors.New(ctx, "no app defined")
	}
	if platformName == "" {
		return errors.New(ctx, "no platform defined")
	}

	eventTypes, err := c.EventTypesList(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to list event types")
	}
	for _, name := range params.SelectedEventNames {
		for _, t := range eventTypes {
			if t.Name == name {
				params.SelectedEventIDs = append(params.SelectedEventIDs, t.ID)
				break
			}
		}
	}

	if len(params.CollaboratorUsernames) > 0 {
		params.UserIDs, err = collaboratorUserIDs(ctx, c, app, params.CollaboratorUsernames)
		if err != nil {
			return errors.Wrapf(ctx, err, "invalid collaborator usernames")
		}
	}

	platforms, err := c.NotificationPlatformByName(ctx, platformName)
	if err != nil {
		return errors.Wrapf(ctx, err, "find notification platform %s", platformName)
	}
	if len(platforms) <= 0 {
		return errors.Newf(ctx, "notification platform \"%s\" has not been found", platformName)
	}
	params.PlatformID = platforms[0].ID

	baseNotifier, err := c.NotifierProvision(ctx, app, params.NotifierParams)
	if err != nil {
		return errors.Wrapf(ctx, err, "create notifier on app %s", app)
	}
	notifier := baseNotifier.Specialize()

	displayDetails(notifier, eventTypes)

	io.Info()
	io.Status("Notifier has been created.")
	return nil
}

func collaboratorUserIDs(ctx context.Context, c *scalingo.Client, app string, usernames []string) ([]string, error) {
	ids := make([]string, 0, len(usernames))

	collaborators, err := c.CollaboratorsList(ctx, app)
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "fail to list collaborators")
	}

	scapp, err := c.AppsShow(ctx, app)
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "fail to get application information")
	}

	var id string
	for _, u := range usernames {
		id = ""
		if u == scapp.Owner.Username {
			id = scapp.Owner.ID
		} else {
			for _, c := range collaborators {
				if c.Username == u && c.Status == "pending" {
					return nil, errors.Newf(ctx, "%v is a collaborator but their invitation is still pending", c.Username)
				} else if c.Username == u {
					id = c.UserID
					break
				}
			}
		}
		if id == "" {
			return nil, errors.Newf(ctx, "no such collaborator: %v", u)
		}

		ids = append(ids, id)
	}

	return ids, nil
}
