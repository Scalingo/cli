package deployments

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/config"
)

func List(app string) error {
	c := config.ScalingoClient()
	deployments, err := c.DeploymentList(app)
	if err != nil {
		return errgo.Mask(err)
	}

	if len(deployments) == 0 {

	} else {
		t := tablewriter.NewWriter(os.Stdout)
		t.SetHeader([]string{"ID", "Date", "User", "Git Ref","Status"})

		for _, deployment := range deployments {
			t.Append([]string{deployment.ID,
				deployment.CreatedAt.Format("2006/01/02 15:04:05"),
				deployment.User.Username,
				deployment.GitRef,
				deployment.Status,
			})
		}
		t.Render()

	}

	return nil
}
