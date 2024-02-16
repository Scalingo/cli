package deployments

import (
	"context"
	stdio "io"
	"os"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-utils/errors/v2"
)

func Logs(ctx context.Context, app, deploymentID string) error {
	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}
	if deploymentID == "" {
		deployments, err := client.DeploymentList(ctx, app)
		if err != nil {
			return errors.Wrapf(ctx, err, "fail to get the most recent deployment")
		}
		if len(deployments) == 0 {
			return errgo.New("This application has not been deployed")
		}
		deploymentID = deployments[0].ID
		io.Infof("-----> Selected the most recent deployment (%s)\n", deploymentID)
	}
	deploy, err := client.Deployment(ctx, app, deploymentID)

	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	res, err := client.DeploymentLogs(ctx, deploy.Links.Output)

	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	defer res.Body.Close()

	if res.StatusCode == 404 {
		io.Error("There is no log for this deployment.")
	} else {
		stdio.Copy(os.Stdout, res.Body)
	}
	return nil
}
