package autocomplete

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func NotifiersAutoComplete(c *cli.Context) error {
	appName := CurrentAppCompletion(c)
	if appName == "" {
		return nil
	}

	client, err := config.ScalingoClient(c.Context)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	resources, err := client.NotifiersList(c.Context, appName)
	if err == nil {
		for _, resource := range resources {
			fmt.Println(resource.GetID())
		}
	}

	return nil
}
