package apps

import (
	"fmt"
	"io"
	"strings"

	"github.com/Scalingo/cli/api"
	"github.com/Scalingo/cli/appdetect"
)

func Create(appName string) error {
	res, err := api.AppsCreate(appName)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 400 {
		err = handleCreateError(appName, res.Body)
		if err != nil {
			return err
		}
	}

	if res.StatusCode == 201 {
		appMap := map[string]App{}
		ReadJson(res.Body, &appMap)
		app := appMap["app"]

		fmt.Printf("App '%s' has been created\n", app.Name)
		if appdetect.DetectGit() && appdetect.AddRemote(app.GitUrl) == nil {
			fmt.Printf("Git repository detected: remote scalingo added\n→ 'git push scalingo master' to deploy your app\n")
		} else {
			fmt.Printf("To deploy your application, run these commands in your GIT repository:\n→ git remote add scalingo %s\n→ git push scalingo master\n", app.GitUrl)
		}
	}
	return nil
}

func handleCreateError(app string, body io.ReadCloser) error {
	errs := api.BadRequestError{}
	err := ReadJson(body, &errs)
	if err != nil {
		return err
	}
	fmt.Printf("Fail to create app %s\n", app)
	for attr, attrErrs := range errs.Errors {
		fmt.Printf("%s → %s\n", attr, strings.Join(attrErrs, ", "))
	}
	return nil
}
