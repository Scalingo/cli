package autocomplete

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-utils/errors/v3"
)

func NotifiersAutoComplete(ctx context.Context, c *cli.Command) error {
	appName := CurrentAppCompletion(c)
	if appName == "" {
		return nil
	}

	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}
	resources, err := client.NotifiersList(ctx, appName)
	if err == nil {
		for _, resource := range resources {
			fmt.Println(resource.GetID())
		}
	}

	return nil
}
