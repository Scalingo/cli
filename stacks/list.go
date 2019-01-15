package stacks

import (
	"os"

	"github.com/Scalingo/cli/config"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"
)

func List() error {
	c := config.ScalingoClient()
	stacks, err := c.StacksList()
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
