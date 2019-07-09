package autocomplete

import (
	"fmt"

	"github.com/Scalingo/cli/config"
	"github.com/urfave/cli"
	"gopkg.in/errgo.v1"
)

func ScaleAutoComplete(c *cli.Context) error {
	appName := CurrentAppCompletion(c)
	if appName == "" {
		return nil
	}

	client, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	processes, err := client.AppsPs(appName)
	if err != nil {
		return errgo.Mask(err)
	}
	for _, ct := range processes {
		fmt.Println(fmt.Sprintf("%s:%d:%s", ct.Name, ct.Amount, ct.Size))
	}

	return nil
}
