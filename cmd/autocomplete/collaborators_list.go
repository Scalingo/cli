package autocomplete

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-utils/errors/v3"
)

func CollaboratorsGenericListAutoComplete(ctx context.Context, c *cli.Command) error {
	appName := CurrentAppCompletion(ctx, c)
	if appName == "" {
		return nil
	}

	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "fail to get Scalingo client")
	}
	collaborators, err := client.CollaboratorsList(ctx, appName)
	if err == nil {
		for _, col := range collaborators {
			fmt.Println(col.Email)
		}
	}

	return nil
}
