package autocomplete

import (
	"context"
	"fmt"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-utils/errors/v3"
)

func StacksSetAutoComplete(ctx context.Context) error {
	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}

	stacks, err := client.StacksList(ctx)
	if err != nil {
		return nil
	}

	for _, stack := range stacks {
		fmt.Println(stack.ID)
		fmt.Println(stack.Name)
	}

	return nil
}
