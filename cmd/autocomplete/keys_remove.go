package autocomplete

import (
	"fmt"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/cli/api"
)

func KeysRemoveAutoComplete(c *cli.Context) error {
	keys, err := api.KeysList()
	if err == nil {

		for _, key := range keys {
			fmt.Println(key.Name)
		}
	}

	return nil
}
