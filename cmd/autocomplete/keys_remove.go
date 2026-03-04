package autocomplete

import (
	"context"
	"fmt"

	"github.com/Scalingo/go-utils/errors/v2"

	"github.com/Scalingo/cli/config"
)

func KeysRemoveAutoComplete(ctx context.Context) error {
	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}
	keys, err := client.KeysList(ctx)
	if err == nil {
		for _, key := range keys {
			fmt.Println(key.Name)
		}
	}

	return nil
}
