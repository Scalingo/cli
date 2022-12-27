package apps

import (
	"context"
	"fmt"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func SendSignal(ctx context.Context, app string, signal string, containerNames []string) error {
	if len(containerNames) == 0 {
		return errgo.New("at least one container name should be given")
	}
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client to send signal to application containers")
	}

	for _, containerName := range containerNames {
		err := c.ContainersSendSignal(ctx, app, signal, containerName)
		if err != nil {
			return errgo.Notef(err, "fail to send signal to container")
		}
		fmt.Printf("-----> Sent signal '%v' to '%v' container.\n", signal, containerName)
	}
	return nil
}
