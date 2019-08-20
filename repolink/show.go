package repolink

import (
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
)

func Show(app string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	repoLink, err := c.SCMRepoLinkShow(app)
	if err != nil {
		return errgo.Notef(err, "fail to get repo link for this app")
	}
	if repoLink == nil {
		io.Statusf("No repo link is linked with '%s' app.\n", app)
		return nil
	}

	i, err := c.SCMIntegrationsShow(repoLink.AuthIntegrationUUID)
	if err != nil {
		return errgo.Notef(err, "fail to get integration of this repo link")
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{
		"App ID", "Integration ID", "Integration Name", "Owner", "Repo", "Branch", "Created At",
		"Auto Deploy", "Review Apps Deploy", "Delete on Close", "Delete on Stale",
	})
	t.Append([]string{
		repoLink.AppID, repoLink.AuthIntegrationUUID, i.SCMType.Str(),
		repoLink.Owner, repoLink.Repo, repoLink.Branch, repoLink.CreatedAt.Format(utils.TimeFormat),
		strconv.FormatBool(repoLink.AutoDeployEnabled), strconv.FormatBool(repoLink.DeployReviewAppsEnabled),
		fmt.Sprintf("%v (%d)", repoLink.DeleteOnCloseEnabled, repoLink.HoursBeforeDeleteOnClose),
		fmt.Sprintf("%v (%d)", repoLink.DeleteStaleEnabled, repoLink.HoursBeforeDeleteStale),
	})
	t.Render()

	return nil
}
