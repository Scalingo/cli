package autocomplete

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/go-utils/errors/v3"
)

func DeploymentsAutoComplete(ctx context.Context, c *cli.Command) error {
	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}
	currentApp := detect.CurrentApp(c)

	deployments, err := client.DeploymentList(ctx, currentApp)
	if err != nil {
		return errors.Wrapf(ctx, err, "list deployments for app %s", currentApp)
	}

	for _, deployment := range deployments {
		fmt.Println(deployment.ID)
	}

	return nil
}
