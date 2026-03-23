package deployments

import (
	"context"
	stdio "io"
	"net/http"
	"os"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	httpclient "github.com/Scalingo/go-scalingo/v11/http"
	"github.com/Scalingo/go-utils/errors/v3"
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
			return errors.New(ctx, "This application has not been deployed")
		}
		deploymentID = deployments[0].ID
		io.Infof("-----> Selected the most recent deployment (%s)\n", deploymentID)
	}

	deploy, err := client.Deployment(ctx, app, deploymentID)
	if err != nil {
		return errors.Wrapf(ctx, err, "get deployment %s", deploymentID)
	}

	readCloser, err := client.DeploymentLogs(ctx, deploy.Links.Output)
	if err != nil {
		var requestFailedErr *httpclient.RequestFailedError
		if !errors.As(err, &requestFailedErr) || requestFailedErr.Code != http.StatusNotFound {
			return errors.Wrap(ctx, err, "fetch deployment logs")
		}

		io.Error("There is no log for this deployment.")
		return nil
	}
	defer readCloser.Close()

	_, err = stdio.Copy(os.Stdout, readCloser)
	if err != nil {
		return errors.Wrap(ctx, err, "copy deployment logs to stdout")
	}
	return nil
}
