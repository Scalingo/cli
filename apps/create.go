package apps

import (
	"fmt"

	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo"
	"gopkg.in/errgo.v1"
)

func Create(appName string, remote string, buildpack string) error {
	c := config.ScalingoClient()
	app, err := c.AppsCreate(scalingo.AppsCreateOpts{Name: appName})
	if err != nil {
		if !utils.PaymentRequiredAndFreeTrialExceeded(err) {
			return errgo.Mask(err, errgo.Any)
		}

		// If error is Payment Required and user tries to exceed its free trial
		validate, err := utils.AskUserValidation()
		if err != nil {
			return errgo.Mask(err, errgo.Any)
		}
		if !validate {
			fmt.Println("Do not break free trial.")
			return nil
		}
		_, err = c.UpdateUser(scalingo.UpdateUserParams{StopFreeTrial: true})
		if err != nil {
			return errgo.Mask(err, errgo.Any)
		}
		return Create(appName, remote, buildpack)
	}

	if buildpack != "" {
		fmt.Println("Installing custom buildpack...")
		_, _, err := c.VariableSet(app.Name, "BUILDPACK_URL", buildpack)
		if err != nil {
			fmt.Println("Failed to set custom buildpack. Please add BUILDPACK_URL=" + buildpack + " to your application environment")
		}
	}

	fmt.Printf("App '%s' has been created\n", app.Name)
	if _, ok := appdetect.DetectGit(); ok && appdetect.AddRemote(app.GitUrl, remote) == nil {
		fmt.Printf("Git repository detected: remote %s added\n→ 'git push %s master' to deploy your app\n", remote, remote)
	} else {
		fmt.Printf("To deploy your application, run these commands in your GIT repository:\n→ git remote add %s %s\n→ git push %s master\n", remote, app.GitUrl, remote)
	}
	return nil
}
