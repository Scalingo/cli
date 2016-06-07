package apps

import (
	"fmt"

	"gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/appdetect"
	"github.com/Scalingo/cli/config"
)

func Create(appName string, remote string) error {
	c := config.ScalingoClient()
	app, err := c.AppsCreate(appName)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	fmt.Printf("App '%s' has been created\n", app.Name)
	if _, ok := appdetect.DetectGit(); ok && appdetect.AddRemote(app.GitUrl, remote) == nil {
		fmt.Printf("Git repository detected: remote %s added\n→ 'git push %s master' to deploy your app\n", remote, remote)
	} else {
		fmt.Printf("To deploy your application, run these commands in your GIT repository:\n→ git remote add %s %s\n→ git push %s master\n", remote, app.GitUrl, remote)
	}
	return nil
}
