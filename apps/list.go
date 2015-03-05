package apps

import (
	"fmt"
	"os"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/olekukonko/tablewriter"
	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/api"
	"github.com/Scalingo/cli/io"
)

func List() error {
	res, err := api.AppsList()
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	defer res.Body.Close()

	appsMap := map[string][]*api.App{}
	err = ReadJson(res.Body, &appsMap)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	apps := appsMap["apps"]

	if len(apps) == 0 {
		fmt.Println(io.Indent("\nYou haven't created any app yet, create your first application using:\n→ scalingo create <app_name>\n", 2))
		return nil
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"Name", "Role", "Owner"})

	for _, app := range apps {
		if app.Owner.Email == api.CurrentUser.Email {
			t.Append([]string{app.Name, "owner", "-"})
		} else {
			t.Append([]string{app.Name, "collaborator", fmt.Sprintf("%s <%s>", app.Owner.Username, app.Owner.Email)})
		}
	}
	t.Render()

	return nil
}
