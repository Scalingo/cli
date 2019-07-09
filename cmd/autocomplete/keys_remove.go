package autocomplete

import (
	"fmt"

	"github.com/Scalingo/cli/config"
	"github.com/urfave/cli"
	"gopkg.in/errgo.v1"
)

func KeysRemoveAutoComplete(c *cli.Context) error {
	client, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	keys, err := client.KeysList()
	if err == nil {
		for _, key := range keys {
			fmt.Println(key.Name)
		}
	}

	return nil
}
