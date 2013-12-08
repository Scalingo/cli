package apps

import (
	"appsdeck/api"
	"appsdeck/appdetect"
	"fmt"
	"strings"
)

func Create(appName string) {
	res, err := api.AppsCreate(appName)
	if err != nil {
		fmt.Println("Fail to create app:", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 400 {
		errs := api.BadRequestError{}
		ReadJson(res.Body, &errs)
		fmt.Printf("Fail to create app %s\n", appName)
		for attr, attrErrs := range errs.Errors {
			fmt.Printf("%s → %s\n", attr, strings.Join(attrErrs, ", "))
		}
	} else if res.StatusCode == 201 {
		app := App{}
		ReadJson(res.Body, &app)
		fmt.Printf("App %s (%s) has been created\n", app.FullName, app.Name)
		if appdetect.DetectGit() && appdetect.AddRemote(app.GitUrl) == nil {
			fmt.Printf("Git repository detected: remote appsdeck added\n→ 'git push appsdeck master' to deploy your app\n")
		} else {
			fmt.Printf("To deploy your application, run these commands in your GIT repository:\n→ git remote add appsdeck %s\n→ git push appsdeck master\n", app.GitUrl)
		}
	}
}
