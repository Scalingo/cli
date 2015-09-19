package autocomplete

import (
	"fmt"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/cli/debug"
)

func FlagAppAutoComplete(c *cli.Context) bool {
	apps, err := appsList()
	if err != nil {
		debug.Println("fail to get apps list:", err)
		return false
	}

	for _, app := range apps {
		fmt.Println(app.Name)
	}

	return true
}
