package scm_integrations

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
)

func List() error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	integrations, err := c.IntegrationsList()
	if err != nil {
		return errgo.Notef(err, "fail to get integrations")
	}

	nbrIntegrations := len(integrations)
	if nbrIntegrations == 0 {
		io.Status("No integration is linked with your Scalingo account.")
		return nil
	}

	io.Statusf("You already have %d integration(s) linked with your Scalingo account:\n", nbrIntegrations)

	t := tablewriter.NewWriter(os.Stdout)
	t.SetColWidth(60)
	t.SetHeader([]string{"ID", "Type", "URL", "Username", "Email"})
	for _, i := range integrations {
		t.Append([]string{i.ID, i.ScmType, i.Url, i.Username, i.Email})
	}
	t.Render()
	return nil
}
