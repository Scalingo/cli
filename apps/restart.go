package apps

import (
	"context"
	"fmt"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-utils/errors/v3"
)

func Restart(ctx context.Context, app string, sync bool, args []string) error {
	params := scalingo.AppsRestartParams{Scope: args}

	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}

	restartOpURL, err := c.AppsRestart(ctx, app, &params)
	if err != nil {
		return errors.Wrapf(ctx, err, "restart app %s", app)
	}

	if !sync {
		fmt.Println("Your application is being restarted.")
		return nil
	}

	waiter := newOperationWaiterFromURL(app, restartOpURL)
	_, err = waiter.WaitOperation(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "wait for restart operation")
	}

	fmt.Println("Your application has been restarted.")
	return nil
}
