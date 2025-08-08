package apps

import (
	"context"
	"fmt"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v8"
)

func Create(ctx context.Context, appName, remote, buildpack, projectID string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	app, err := c.AppsCreate(ctx, scalingo.AppsCreateOpts{Name: appName, ProjectID: projectID})
	if err != nil {
		if utils.IsRegionDisabledError(err) {
			return handleRegionDisabledError(ctx, appName, c)
		}
		if !utils.IsPaymentRequiredAndFreeTrialExceededError(err) {
			return errgo.Notef(err, "fail to create the application")
		}
		// If error is Payment Required and user tries to exceed its free trial
		return utils.AskAndStopFreeTrial(ctx, c, func() error {
			return Create(ctx, appName, remote, buildpack, projectID)
		})
	}

	if buildpack != "" {
		fmt.Println("Installing custom buildpack...")
		_, _, err := c.VariableSet(ctx, app.Name, "BUILDPACK_URL", buildpack)
		if err != nil {
			fmt.Println("Failed to set custom buildpack. Please add BUILDPACK_URL=" + buildpack + " to your application environment")
		}
	}

	fmt.Printf("App '%s' has been created\n", app.Name)
	if _, ok := utils.DetectGit(); ok && utils.AddGitRemote(app.GitURL, remote) == nil {
		fmt.Printf("Git repository detected: remote %s added\n→ 'git push %s master' to deploy your app\n", remote, remote)
	} else {
		fmt.Printf("To deploy your application, run these commands in your GIT repository:\n→ git remote add %s %s\n→ git push %s master\n", remote, app.GitURL, remote)
	}
	return nil
}

func handleRegionDisabledError(ctx context.Context, appName string, c *scalingo.Client) error {
	regions, rerr := c.RegionsList(ctx)
	if rerr != nil {
		return errgo.Notef(rerr, "region is disabled, failed to list available regions")
	}
	if len(regions) <= 1 {
		return errgo.New("region is disabled and there is no other available region")
	}
	firstRegion := regions[0]
	if firstRegion.Name == config.C.ScalingoRegion {
		firstRegion = regions[1]
	}

	fmt.Printf("Application creation has been disabled on the currently used region: %v\n\n", config.C.ScalingoRegion)
	fmt.Printf("Either configure your CLI to use another default region, then create your application:\n")
	fmt.Printf("    scalingo config --region %s\n    scalingo create %s\n", firstRegion.Name, appName)
	fmt.Printf("\nOr use the region flag to specify the region explicitly for this command:\n")
	fmt.Printf("    scalingo --region %s create %s\n", firstRegion.Name, appName)
	fmt.Printf("\nList of available regions:\n")
	for _, region := range regions {
		if region.Name == config.C.ScalingoRegion {
			continue
		}
		fmt.Printf("- %v (%v)\n", region.Name, region.DisplayName)
	}
	return nil
}
