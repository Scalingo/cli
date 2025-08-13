package autocomplete

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func EnvUnsetAutoComplete(ctx context.Context, c *cli.Command) error {
	appName := CurrentAppCompletion(c)
	if appName == "" {
		return nil
	}

	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	variables, err := client.VariablesList(ctx, appName)
	if err == nil {
		for _, v := range variables {
			fmt.Println(v.Name)
		}
	}

	return nil
}
