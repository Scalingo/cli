package autocomplete

import (
	"context"
	"fmt"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-utils/errors/v2"
)

func DatabasesNgListAutoComplete(ctx context.Context) error {
	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	databases, err := client.Preview().DatabasesList(ctx)
	if err == nil {
		for _, db := range databases {
			fmt.Println(db.DatabaseInfo.ID)
		}
	}

	return nil
}
