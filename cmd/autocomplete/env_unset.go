package autocomplete

import (
	"fmt"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/cli/config"
)

func EnvUnsetAutoComplete(c *cli.Context) error {
	appName := CurrentAppCompletion(c)
	if appName == "" {
		return nil
	}

	client := config.ScalingoClient()
	variables, err := client.VariablesList(appName)
	if err == nil {

		for _, v := range variables {
			fmt.Println(v.Name)
		}
	}

	return nil
}
