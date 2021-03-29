package apps

import (
	"fmt"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v4"
	"gopkg.in/errgo.v1"
)

func OneOffStop(appName, oneOffLabel string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client to stop a running one-off")
	}

	containers, err := c.AppsContainersPs(appName)
	if err != nil {
		return errgo.Notef(err, "fail to list the application containers to get the ID of the container to stop")
	}

	var containerToStop *scalingo.Container
	for _, container := range containers {
		if container.Label == oneOffLabel {
			containerToStop = &container
			break
		}
	}
	if containerToStop == nil {
		return fmt.Errorf("The container '%s' does not exist", oneOffLabel)
	}

	err = c.ContainersStop(appName, containerToStop.ID)
	if err != nil {
		// TODO filter if it's a 400 error, do not add a Notef.
		return errgo.Notef(err, "fail to stop the container '%s'", oneOffLabel)
	}

	io.Statusf("Container one-off %v of the app %v is being asynchronously stopped...\n", io.Bold(containerToStop.Label), io.Bold(appName))

	return nil
}
