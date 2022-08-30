package apps

import (
	"context"
	"fmt"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v4"
)

func OneOffStop(ctx context.Context, appName, oneOffLabel string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client to stop a running one-off")
	}

	containers, err := c.AppsContainersPs(ctx, appName)
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

	err = c.ContainersStop(ctx, appName, containerToStop.ID)
	if err != nil {
		return errgo.Notef(err, "fail to stop the container '%s'", oneOffLabel)
	}

	io.Statusf("Container one-off %v of the app %v is being asynchronously stopped...\n", io.Bold(containerToStop.Label), io.Bold(appName))

	return nil
}
