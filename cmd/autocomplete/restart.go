package autocomplete

import (
	"fmt"

	"github.com/urfave/cli"
	"gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/config"
)

func RestartAutoComplete(c *cli.Context) error {
	appName := CurrentAppCompletion(c)
	if appName == "" {
		return nil
	}

	client := config.ScalingoClient()
	processes, err := client.AppsPs(appName)
	if err != nil {
		return errgo.Mask(err)
	}
	for _, ct := range processes {
		fmt.Println(ct.Name)
	}

	return nil
}
