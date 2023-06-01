package autocomplete

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func ScaleAutoComplete(c *cli.Context) error {
	appName := CurrentAppCompletion(c)
	if appName == "" {
		return nil
	}

	client, err := config.ScalingoClient(c.Context)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	processes, err := client.AppsContainerTypes(c.Context, appName)
	if err != nil {
		return errgo.Mask(err)
	}
	for _, ct := range processes {
		fmt.Printf("%s:%d:%s\n", ct.Name, ct.Amount, ct.Size)
	}

	return nil
}
