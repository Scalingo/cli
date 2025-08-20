package autocomplete

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/detect"
)

func DeploymentsAutoComplete(ctx context.Context, c *cli.Command) error {
	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	currentApp := detect.CurrentApp(c)

	deployments, err := client.DeploymentList(ctx, currentApp)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	for _, deployment := range deployments {
		fmt.Println(deployment.ID)
	}

	return nil
}
