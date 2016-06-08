package autocomplete

import (
	"fmt"

	"github.com/Scalingo/codegangsta-cli"
	"gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/config"
  "github.com/Scalingo/cli/appdetect"
)

func DeploymentsAutoComplete(c *cli.Context) error {
  client := config.ScalingoClient()
  currentApp := appdetect.CurrentApp(c)

  deployments, err := client.DeploymentList(currentApp)
  if err != nil {
    return errgo.Mask(err, errgo.Any)
  }

  for _, deployment := range(deployments) {
    fmt.Println(deployment.ID)
  }

  return nil
}
