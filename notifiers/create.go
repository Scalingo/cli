package notifiers

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo/debug"
	"github.com/Scalingo/cli/io"
	scalingo "github.com/Scalingo/go-scalingo"
	"gopkg.in/errgo.v1"
)

type ProvisionParams struct {
	CollaboratorUsernames []string
	SelectedEventNames    []string
	scalingo.NotifierParams
}

func Provision(app, platformName string, params ProvisionParams) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	debug.Printf("[Provision] params: %+v", params)

	if app == "" {
		return errgo.New("no app defined")
	}
	if platformName == "" {
		return errgo.New("no platform defined")
	}

	eventTypes, err := c.EventTypesList()
	if err != nil {
		return errgo.Notef(err, "fail to list event types")
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
		params.UserIDs, err = collaboratorUserIDs(c, app, params.CollaboratorUsernames)
		if err != nil {
			return errgo.Notef(err, "invalid collaborator usernames")
		}
	}

	platforms, err := c.NotificationPlatformByName(platformName)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	if len(platforms) <= 0 {
		return errgo.Newf("notification platform \"%s\" has not been found", platformName)
	}
	params.PlatformID = platforms[0].ID

	baseNotifier, err := c.NotifierProvision(app, params.NotifierParams)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	notifier := baseNotifier.Specialize()

	displayDetails(notifier, eventTypes)

	io.Info()
	io.Status("Notifier have been created.")
	return nil
}

func collaboratorUserIDs(c *scalingo.Client, app string, usernames []string) ([]string, error) {
	ids := make([]string, 0, len(usernames))

	collaborators, err := c.CollaboratorsList(app)
	if err != nil {
		return nil, errgo.Notef(err, "fail to list collaborators")
	}

	scapp, err := c.AppsShow(app)
	if err != nil {
		return nil, errgo.Notef(err, "fail to get application information")
	}

	var id string
	for _, u := range usernames {
		id = ""
		if u == scapp.Owner.Username {
			id = scapp.Owner.ID
		} else {
			for _, c := range collaborators {
				if c.Username == u && c.Status == "pending" {
					return nil, errgo.Newf("%v is a collaborator but their invitation is still pending", c.Username)
				} else if c.Username == u {
					id = c.UserID
					break
				}
			}
		}
		if id == "" {
			return nil, errgo.Newf("no such collaborator: %v", u)
		} else {
			ids = append(ids, id)
		}
	}

	return ids, nil
}
