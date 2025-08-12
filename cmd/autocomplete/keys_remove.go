package autocomplete

import (
	"fmt"

	"github.com/urfave/cli/v3"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func KeysRemoveAutoComplete(c *cli.Context) error {
	client, err := config.ScalingoClient(c.Context)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	keys, err := client.KeysList(c.Context)
	if err == nil {
		for _, key := range keys {
			fmt.Println(key.Name)
		}
	}

	return nil
}
