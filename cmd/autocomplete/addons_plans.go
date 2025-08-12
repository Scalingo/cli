package autocomplete

import (
	"fmt"

	"github.com/urfave/cli/v3"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func AddonsPlansAutoComplete(c *cli.Context) error {
	appName := CurrentAppCompletion(c)
	if appName == "" {
		return nil
	}

	client, err := config.ScalingoClient(c.Context)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	resources, err := client.AddonsList(c.Context, appName)
	if err == nil {
		for _, resource := range resources {
			fmt.Println(resource.AddonProvider.ID)
		}
	}

	return nil
}
