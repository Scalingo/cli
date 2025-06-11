package scmintegrations

import (
	"context"
	"os"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v7"
)

func List(ctx context.Context) error {
	c, err := config.ScalingoAuthClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	integrations, err := c.SCMIntegrationsList(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to list SCM integrations")
	}

	nbrIntegrations := len(integrations)
	if nbrIntegrations == 0 {
		io.Status("Your Scalingo account is not linked to any SCM integrations.")
		return nil
	}

	pluralIntegration := ""
	if nbrIntegrations > 1 {
		pluralIntegration = "s"
	}

	io.Statusf("You already have %d SCM integration%s linked with your Scalingo account:\n", nbrIntegrations, pluralIntegration)

	t := tablewriter.NewWriter(os.Stdout)
	t.Header([]string{"ID", "Type", "URL", "Username", "Email"})
	for _, i := range integrations {
		t.Append([]string{i.ID, scalingo.SCMTypeDisplay[i.SCMType], i.URL, i.Username, i.Email})
	}
	t.Render()
	return nil
}
