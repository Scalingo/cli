package autocomplete

import (
	"fmt"

	"github.com/urfave/cli"

	"github.com/Scalingo/go-scalingo/v4/debug"
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
