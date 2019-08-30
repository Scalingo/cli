package review_apps

import (
	"fmt"
	"os"

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

	reviewApps, err := c.SCMRepoLinkReviewApps(app)
	if err != nil {
		return errgo.Notef(err, "fail to get review apps for this app")
	}
	if len(reviewApps) == 0 {
		io.Statusf("No review app for '%s' or specified app is not a parent app.\n", app)
		return nil
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"App", "PR", "PR Branch", "Created At", "Status"})
	for _, ra := range reviewApps {
		date := ra.CreatedAt.Local().Format(utils.TimeFormat)

		t.Append([]string{
			ra.AppName, fmt.Sprintf("%d", ra.PullRequest.Number), ra.PullRequest.BranchName,
			date, fmt.Sprintf("%v", ra.LastDeployment.Status),
		})
	}
	t.Render()

	return nil
}
