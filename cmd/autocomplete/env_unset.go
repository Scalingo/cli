package autocomplete

import (
	"fmt"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/cli/api"
)

func EnvUnsetAutoComplete(c *cli.Context) error {
	appName := CurrentAppCompletion(c)
	if appName == "" {
		return nil
	}

	variables, err := api.VariablesList(appName)
	if err == nil {

		for _, v := range variables {
			fmt.Println(v.Name)
		}
	}

	return nil
}
