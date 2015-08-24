package autocomplete

import (
	"fmt"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/cli/api"
)

func FlagAppAutoComplete(c *cli.Context) bool {
	apps, err := api.AppsList()
	if err != nil || len(apps) == 0 {
		return false
	}

	for _, app := range apps {
		fmt.Println(app.Name)
	}

	return true
}
