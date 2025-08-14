package autocomplete

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func LogDrainsRemoveAutoComplete(ctx context.Context, c *cli.Command) error {
	appName := CurrentAppCompletion(c)
	if appName == "" {
		return nil
	}

	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	drains, err := client.LogDrainsList(ctx, appName)
	if err != nil {
		return errgo.Notef(err, "fail to get log drains list")
	}

	for _, drain := range drains {
		fmt.Println(drain.URL)
	}

	return nil
}
