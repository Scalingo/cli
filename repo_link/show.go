package repo_link

import (
	"fmt"
	"os"
	"strconv"
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

	repoLink, err := c.ScmRepoLinkShow(app)
	if err != nil {
		return errgo.Notef(err, "fail to get repo link for this app")
	}
	if repoLink == nil {
		fmt.Printf("No repo link is linked with '%s' app.\n", app)
		return nil
	}

	i, err := c.IntegrationsShow(repoLink.AuthIntegrationID)
	if err != nil {
		return errgo.Notef(err, "fail to get integration of this repo link")
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{
		"App ID", "Integration ID", "Integration Name", "Owner", "Repo", "Branch", "Created At",
		"Auto Deploy", "Review Apps Deploy", "Delete on close", "Delete on stale",
	})
	t.Append([]string{
		repoLink.AppID, repoLink.AuthIntegrationID, i.ScmType,
		repoLink.Owner, repoLink.Repo, repoLink.Branch, repoLink.CreatedAt.Format(time.RFC1123),
		strconv.FormatBool(repoLink.AutoDeployEnabled), strconv.FormatBool(repoLink.DeployReviewAppsEnabled),
		fmt.Sprintf("%v (%d)", repoLink.DeleteOnCloseEnabled, repoLink.HoursBeforeDeleteOnClose),
		fmt.Sprintf("%v (%d)", repoLink.DeleteStaleEnabled, repoLink.HoursBeforeDeleteStale),
	})
	t.Render()

	return nil
}
