package notifiers

import (
	"context"
	"fmt"
	"strings"

	"github.com/Scalingo/go-utils/errors/v3"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v10"
)

func Destroy(ctx context.Context, app string, ID string) error {
	if app == "" {
		return errors.New(ctx, "no app defined")
	}
	if ID == "" {
		return errors.New(ctx, "no ID defined")
	}
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}

	notifier, err := getNotifier(ctx, c, app, ID)
	if err != nil {
		return errors.Wrapf(ctx, err, "find notifier %s on app %s", ID, app)
	}

	io.Status("Destroy", notifier.GetName())
	io.Warning("This operation is irreversible")
	io.Warning("All related data will be destroyed")
	io.Info("To confirm, type the ID of the notifier (" + ID + "):")
	fmt.Print("-----> ")

	var validationID string
	fmt.Scan(&validationID)

	if validationID != ID {
		return errors.Newf(ctx, "'%s' is not '%s', aborting…\n", validationID, ID)
	}

	err = c.NotifierDestroy(ctx, app, notifier.GetID())
	if err != nil {
		return errors.Wrapf(ctx, err, "delete notifier %s on app %s", ID, app)
	}

	io.Status("Notifier", ID, "has been destroyed")
	return nil
}

func getNotifier(ctx context.Context, c *scalingo.Client, app string, ID string) (scalingo.DetailedNotifier, error) {
	resources, err := c.NotifiersList(ctx, app)
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "list notifiers on app %s", app)
	}
	notifiersList := []string{}
	for _, r := range resources {
		notifiersList = append(notifiersList, fmt.Sprintf("%s: [%s] %s", r.GetID(), string(r.GetType()), r.GetName()))
		if ID == r.GetID() {
			return r, nil
		}
	}
	return nil, errors.Newf(ctx, "Notifier %s doesn't exist for app %s\nExisting notifiers:\n  - %v", ID, app, strings.Join(notifiersList, "\n  - "))
}
