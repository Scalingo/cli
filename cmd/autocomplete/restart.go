package autocomplete

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-utils/errors/v3"
)

func RestartAutoComplete(ctx context.Context, c *cli.Command) error {
	appName := CurrentAppCompletion(ctx, c)
	if appName == "" {
		return nil
	}

	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}
	processes, err := client.AppsContainerTypes(ctx, appName)
	if err != nil {
		return errors.Wrapf(ctx, err, "list container types for app %s", appName)
	}
	for _, ct := range processes {
		fmt.Println(ct.Name)
	}

	return nil
}
