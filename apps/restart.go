package apps

import (
	"context"
	"fmt"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo/v8"
)

func Restart(ctx context.Context, app string, sync bool, args []string) error {
	params := scalingo.AppsRestartParams{Scope: args}

	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	res, err := c.AppsRestart(ctx, app, &params)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	res.Body.Close()

	if !sync {
		fmt.Println("Your application is being restarted.")
		return nil
	}

	waiter := NewOperationWaiterFromHTTPResponse(app, res)
	_, err = waiter.WaitOperation(ctx)
	if err != nil {
		return errgo.Mask(err)
	}

	fmt.Println("Your application has been restarted.")
	return nil
}
