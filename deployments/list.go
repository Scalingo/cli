package deployments

import (
	"fmt"
	"os"
	"time"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/utils"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"
)

func List(app string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	deployments, err := c.DeploymentList(app)
	if err != nil {
		return errgo.Mask(err)
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"ID", "Date", "Duration", "User", "Git Ref", "Status"})

	for _, deployment := range deployments {
		var duration string
		if deployment.Duration != 0 {
			d, _ := time.ParseDuration(fmt.Sprintf("%ds", deployment.Duration))
			duration = d.String()
		} else {
			duration = "n/a"
		}
		t.Append([]string{deployment.ID,
			deployment.CreatedAt.Format(utils.TimeFormat),
			duration,
			deployment.User.Username,
			deployment.GitRef,
			string(deployment.Status),
		})
	}
	t.Render()

	return nil
}
