package repo_link

import (
	"fmt"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func Show(app string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	rl, err := c.ScmRepoLinkShow(app)
	if err != nil {
		return errgo.Mask(err)
	}

	if rl == nil {
		fmt.Printf("No repo link is linked with '%s' app.\n", app)
		return nil
	}

	i, err := c.IntegrationsShow(rl.AuthIntegrationID)
	if err != nil {
		return errgo.Mask(err)
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{
		"ID", "App ID", "Integration ID", "Integration Name", "Owner", "Repo", "Branch", "Created At",
	})
	t.Append([]string{
		rl.ID, rl.AppID, rl.AuthIntegrationID, i.ScmType,
		rl.Owner, rl.Repo, rl.Branch, rl.CreatedAt.Format(time.RFC1123),
	})
	t.Render()

	return nil
}
