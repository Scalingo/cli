package autocomplete

import (
	"context"
	"fmt"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func StacksSetAutoComplete(ctx context.Context) error {
	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
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
