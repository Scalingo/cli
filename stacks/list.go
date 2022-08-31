package stacks

import (
	"context"
	"os"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func List(ctx context.Context) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	stacks, err := c.StacksList(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to list available stacks")
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"ID", "Name", "Description", "Default?"})

	for _, stack := range stacks {
		defaultText := "No"
		if stack.Default {
			defaultText = "Yes"
		}
		t.Append([]string{stack.ID, stack.Name, stack.Description, defaultText})
	}
	t.Render()
	return nil
}
