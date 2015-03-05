package apps

import (
	"fmt"

	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/api"
	"github.com/Scalingo/cli/appdetect"
)

func Create(appName string) error {
	app, err := api.AppsCreate(appName)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	fmt.Printf("App '%s' has been created\n", app.Name)
	if appdetect.DetectGit() && appdetect.AddRemote(app.GitUrl) == nil {
		fmt.Printf("Git repository detected: remote scalingo added\n→ 'git push scalingo master' to deploy your app\n")
	} else {
		fmt.Printf("To deploy your application, run these commands in your GIT repository:\n→ git remote add scalingo %s\n→ git push scalingo master\n", app.GitUrl)
	}
	return nil
}
