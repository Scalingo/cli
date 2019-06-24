package autocomplete

import (
	"fmt"

	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/config"
	"github.com/urfave/cli"
	"gopkg.in/errgo.v1"
)

func DeploymentsAutoComplete(c *cli.Context) error {
	client, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	currentApp := appdetect.CurrentApp(c)

	deployments, err := client.DeploymentList(currentApp)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	for _, deployment := range deployments {
		fmt.Println(deployment.ID)
	}

	return nil
}
