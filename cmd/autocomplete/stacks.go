package autocomplete

import (
	"fmt"

	"github.com/urfave/cli/v3"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func StacksSetAutoComplete(c *cli.Context) error {
	client, err := config.ScalingoClient(c.Context)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	stacks, err := client.StacksList(c.Context)
	if err != nil {
		return nil
	}

	for _, stack := range stacks {
		fmt.Println(stack.ID)
		fmt.Println(stack.Name)
	}

	return nil
}
