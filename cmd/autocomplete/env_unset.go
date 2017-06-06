package autocomplete

import (
	"fmt"

	"github.com/urfave/cli"
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
