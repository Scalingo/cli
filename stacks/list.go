package stacks

import (
	"context"
	"os"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func List(ctx context.Context, withDeprecated bool) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	stacks, err := c.StacksList(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to list available stacks")
	}

	t := tablewriter.NewWriter(os.Stdout)

	if withDeprecated {
		t.SetHeader([]string{"ID", "Name", "Description", "Default?", "Deprecated?", "Deprecation date"})
	} else {
		t.SetHeader([]string{"ID", "Name", "Description", "Default?"})
	}

	for _, stack := range stacks {
		defaultText := "No"
		if stack.Default {
			defaultText = "Yes"
		}

		if withDeprecated {
			deprecatedText := "No"
			deprecationDate := ""

			if stack.IsDeprecated() {
				deprecatedText = "Yes"
			}

			if !stack.DeprecatedAt.IsZero() {
				deprecationDate = stack.DeprecatedAt.Format("2006-01-02")
			}

			t.Append([]string{stack.ID, stack.Name, stack.Description, defaultText, deprecatedText, deprecationDate})
		} else {
			t.Append([]string{stack.ID, stack.Name, stack.Description, defaultText})
		}
	}
	t.Render()
	return nil
}
