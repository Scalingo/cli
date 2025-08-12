package autocomplete

import (
	"fmt"

	"github.com/urfave/cli/v3"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func EnvUnsetAutoComplete(c *cli.Context) error {
	appName := CurrentAppCompletion(c)
	if appName == "" {
		return nil
	}

	client, err := config.ScalingoClient(c.Context)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	variables, err := client.VariablesList(c.Context, appName)
	if err == nil {
		for _, v := range variables {
			fmt.Println(v.Name)
		}
	}

	return nil
}
