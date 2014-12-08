package apps

import (
	"fmt"
	"os"

	"github.com/Scalingo/cli/api"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"
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
