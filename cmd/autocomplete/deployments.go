package autocomplete

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/detect"
)

func DeploymentsAutoComplete(c *cli.Context) error {
	client, err := config.ScalingoClient(c.Context)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	currentApp := detect.CurrentApp(c)

	deployments, err := client.DeploymentList(c.Context, currentApp)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	for _, deployment := range deployments {
		fmt.Println(deployment.ID)
	}

	return nil
}
