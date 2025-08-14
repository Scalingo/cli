package autocomplete

import (
	"context"
	"fmt"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func KeysRemoveAutoComplete(ctx context.Context) error {
	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	keys, err := client.KeysList(ctx)
	if err == nil {
		for _, key := range keys {
			fmt.Println(key.Name)
		}
	}

	return nil
}
