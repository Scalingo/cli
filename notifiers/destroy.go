package notifiers

import (
	"fmt"
	"strings"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v4"
	"gopkg.in/errgo.v1"
)

func Destroy(app string, ID string) error {
	if app == "" {
		return errgo.New("no app defined")
	}
	if ID == "" {
		return errgo.New("no ID defined")
	}
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	notifier, err := getNotifier(c, app, ID)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	io.Status("Destroy", notifier.GetName())
	io.Warning("This operation is irreversible")
	io.Warning("All related data will be destroyed")
	io.Info("To confirm, type the ID of the notifier (" + ID + "):")
	fmt.Print("-----> ")

	var validationID string
	fmt.Scan(&validationID)

	if validationID != ID {
		return errgo.Newf("'%s' is not '%s', aborting…\n", validationID, ID)
	}

	err = c.NotifierDestroy(app, notifier.GetID())
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	io.Status("Notifier", ID, "has been destroyed")
	return nil
}

func getNotifier(c *scalingo.Client, app string, ID string) (scalingo.DetailedNotifier, error) {
	resources, err := c.NotifiersList(app)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	notifiersList := []string{}
	for _, r := range resources {
		notifiersList = append(notifiersList, fmt.Sprintf("%s: [%s] %s", r.GetID(), string(r.GetType()), r.GetName()))
		if ID == r.GetID() {
			return r, nil
		}
	}
	return nil, errgo.Newf("Notifier %s doesn't exist for app %s\nExisting notifiers:\n  - %v", ID, app, strings.Join(notifiersList, "\n  - "))
}
