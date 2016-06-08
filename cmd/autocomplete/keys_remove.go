package autocomplete

import (
	"fmt"

	"github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/cli/config"
)

func KeysRemoveAutoComplete(c *cli.Context) error {
	client := config.ScalingoClient()
	keys, err := client.KeysList()
	if err == nil {
		for _, key := range keys {
			fmt.Println(key.Name)
		}
	}

	return nil
}
