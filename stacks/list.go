package stacks

import (
	"os"
	"strings"

	"github.com/Scalingo/cli/config"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"
)

func List() error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	stacks, err := c.StacksList()
	if err != nil {
		return errgo.Notef(err, "fail to list available stacks")
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"ID", "Name", "Description", "Aliases", "Default?"})

	for _, stack := range stacks {
		defaultText := "No"
		if stack.Default {
			defaultText = "Yes"
		}
		t.Append([]string{stack.ID, stack.Name, stack.Description, strings.Join(stack.Aliases, ", "), defaultText})
	}
	t.Render()
	return nil
}
