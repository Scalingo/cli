package notifications

import (
	"fmt"
	"strings"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/go-scalingo"
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/io"
)

func Destroy(app string, ID string) error {
	if app == "" {
		return errgo.New("no app defined")
	} else if ID == "" {
		return errgo.New("no ID defined")
	}

	notification, err := checkNotificationExist(app, ID)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	io.Status("Destroy", notification.WebHookURL)
	io.Warning("This operation is irreversible")
	io.Warning("All related data will be destroyed")
	io.Info("To confirm, type the ID of the notification (" + ID + "):")
	fmt.Print("-----> ")

	var validationID string
	fmt.Scan(&validationID)

	if validationID != ID {
		return errgo.Newf("'%s' is not '%s', aborting…\n", validationID, ID)
	}

	c := config.ScalingoClient()
	err = c.NotificationDestroy(app, notification.ID)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	io.Status("Notification", ID, "has been destroyed")
	return nil
}

func checkNotificationExist(app string, ID string) (*scalingo.Notification, error) {
	c := config.ScalingoClient()
	resources, err := c.NotificationsList(app)
	if err != nil {
		return nil, errgo.Mask(err, errgo.Any)
	}
	notificationsList := []string{}
	for _, r := range resources {
		notificationsList = append(notificationsList, r.Type + " (" + r.WebHookURL + ")")
		if ID == r.ID {
			return r, nil
		}
	}
	return nil, errgo.Newf("Notification " + ID + " doesn't exist for app " + app + "\nExisting notifications:\n  - %v", strings.Join(notificationsList, "\n  - "))
}
