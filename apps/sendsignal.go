package apps

import (
	"context"
	"fmt"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func SendSignal(ctx context.Context, app string, signal string, args []string) error {
	if len(args) == 0 {
		return errgo.New("at least one container name should be given")
	}
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client to send signal to application containers")
	}

	for _, container := range args {
		err := c.AppsContainersSendSignal(ctx, app, signal, container)
		if err != nil {
			return errgo.Notef(err, "fail to send signal to container")
		}
		fmt.Printf("-----> Sending signal '%v' to '%v' container.\n", signal, container)
	}
	return nil
}
