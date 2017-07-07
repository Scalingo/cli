package autocomplete

import (
	"fmt"

	"github.com/urfave/cli"
	"github.com/Scalingo/cli/config"
)

func NotificationsRemoveAutoComplete(c *cli.Context) error {
	appName := CurrentAppCompletion(c)
	if appName == "" {
		return nil
	}

	client := config.ScalingoClient()
	resources, err := client.NotificationsList(appName)
	if err == nil {
		for _, resource := range resources {
			fmt.Println(resource.ID)
		}
	}

	return nil
}
