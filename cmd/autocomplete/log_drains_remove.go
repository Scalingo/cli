package autocomplete

import (
	"fmt"

	"github.com/urfave/cli"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func LogDrainsRemoveAutoComplete(c *cli.Context) error {
	appName := CurrentAppCompletion(c)
	if appName == "" {
		return nil
	}

	client, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	drains, err := client.LogDrainsList(appName)
	if err != nil {
		return errgo.Notef(err, "fail to get log drains list")
	}

	for _, drain := range drains {
		fmt.Println(drain.URL)
	}

	return nil
}
