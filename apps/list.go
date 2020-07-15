package apps

import (
	"fmt"
	"os"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"
)

func List() error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	apps, err := c.AppsList()
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	if len(apps) == 0 {
		fmt.Println(io.Indent("\nYou haven't created any app yet, create your first application using:\n→ scalingo create <app_name>\n", 2))
		return nil
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"Name", "Role", "Status"})

	currentUser, err := config.C.CurrentUser()
	if err != nil {
		return errgo.Notef(err, "fail to get current user")
	}

	for _, app := range apps {
		if app.Owner.Email == currentUser.Email {
			t.Append([]string{app.Name, "owner", string(app.Status)})
		} else {
			t.Append([]string{app.Name, "collaborator", string(app.Status)})
		}
	}
	t.Render()

	return nil
}
