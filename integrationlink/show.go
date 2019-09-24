package integrationlink

import (
	"fmt"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
	scalingo "github.com/Scalingo/go-scalingo"
	"github.com/fatih/color"
)

func Show(app string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	repoLink, err := c.SCMRepoLinkShow(app)
	if err != nil {
		return errgo.Notef(err, "fail to get integration link for this app")
	}
	if repoLink == nil {
		io.Statusf("Your app '%s' has no integration link.\n", app)
		return nil
	}

	fmt.Printf("%s: %s (%s)\n",
		color.New(color.FgYellow).Sprint("Application"),
		app, repoLink.AppID,
	)
	fmt.Printf("%s: %s (%s)\n",
		color.New(color.FgYellow).Sprint("Integration"),
		scalingo.SCMTypeDisplay[repoLink.SCMType], repoLink.AuthIntegrationUUID,
	)
	fmt.Println()

	fmt.Printf("%s: %s/%s\n",
		color.New(color.FgYellow).Sprint("Repository"),
		repoLink.Owner, repoLink.Repo,
	)
	var autoDeploy string
	if repoLink.AutoDeployEnabled {
		autoDeploy = fmt.Sprintf("%s %s", color.GreenString(utils.Success), repoLink.Branch)
	} else {
		autoDeploy = color.RedString(utils.Error)
	}
	fmt.Printf("%s: %s\n",
		color.New(color.FgYellow).Sprint("Auto Deploy"),
		autoDeploy,
	)

	var reviewAppsDeploy string
	if repoLink.DeployReviewAppsEnabled {
		reviewAppsDeploy = color.GreenString(utils.Success)
	} else {
		reviewAppsDeploy = color.RedString(utils.Error)
	}
	fmt.Printf("%s: %v\n",
		color.New(color.FgYellow).Sprint("Review Apps Deploy"),
		reviewAppsDeploy,
	)
	if repoLink.DeployReviewAppsEnabled {
		var deleteOnClose string
		if repoLink.DeleteOnCloseEnabled {
			if repoLink.HoursBeforeDeleteOnClose == 0 {
				deleteOnClose = "instantly"
			} else {
				deleteOnClose = fmt.Sprintf("after %dh", repoLink.HoursBeforeDeleteOnClose)
			}
		} else {
			deleteOnClose = color.RedString(utils.Error)
		}
		fmt.Printf("\t%s: %s\n",
			color.New(color.FgYellow).Sprint("Destroy on Close"),
			deleteOnClose,
		)

		var deleteOnStale string
		if repoLink.DeleteStaleEnabled {
			if repoLink.HoursBeforeDeleteStale == 0 {
				deleteOnStale = "instantly"
			} else {
				deleteOnStale = fmt.Sprintf("after %dh", repoLink.HoursBeforeDeleteStale)
			}
		} else {
			deleteOnStale = color.RedString(utils.Error)
		}
		fmt.Printf("\t%s: %s\n",
			color.New(color.FgYellow).Sprint("Destroy on Stale"),
			deleteOnStale,
		)
	}

	return nil
}
