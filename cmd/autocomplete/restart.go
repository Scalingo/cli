package autocomplete

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/go-utils/errors/v2"

	"github.com/Scalingo/cli/config"
)

func RestartAutoComplete(ctx context.Context, c *cli.Command) error {
	appName := CurrentAppCompletion(c)
	if appName == "" {
		return nil
	}

	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}
	processes, err := client.AppsContainerTypes(ctx, appName)
	if err != nil {
		return errors.Wrap(ctx, err, "operation failed")
	}
	for _, ct := range processes {
		fmt.Println(ct.Name)
	}

	return nil
}
