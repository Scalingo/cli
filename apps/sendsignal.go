package apps

import (
	"context"
	"fmt"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo/v6"
)

func SendSignal(ctx context.Context, appName string, signal string, containerNames []string) error {
	if len(containerNames) == 0 {
		return errgo.New("at least one container name should be given")
	}
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client to send signal to application containers")
	}

	containers, err := c.AppsContainersPs(ctx, appName)
	if err != nil {
		return errgo.Notef(err, "fail to list the application containers to get the ID of the container to stop")
	}

	for _, containerName := range containerNames {
		var containerToStop *scalingo.Container
		for _, container := range containers {
			if container.Label == containerName {
				containerToStop = &container
				break
			}
		}
		if containerToStop == nil {
			return fmt.Errorf("The container '%s' does not exist", containerName)
		}

		err := c.ContainersSendSignal(ctx, appName, signal, containerToStop.ID)
		if err != nil {
			return errgo.Notef(err, "fail to send signal to container")
		}
		fmt.Printf("-----> Sent signal '%v' to '%v' container.\n", signal, containerName)
	}
	return nil
}
