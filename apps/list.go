package apps

import (
	"github.com/Appsdeck/appsdeck/api"
	"github.com/Appsdeck/appsdeck/auth"
	"fmt"
)

func List() error {
	res, _ := api.AppsList()
	defer res.Body.Close()

	apps := []App{}
	ReadJson(res.Body, &apps)

	fmt.Printf("List of your apps :\n")
	for _, app := range apps {
		if app.Owner.Email == auth.Config.Email {
			fmt.Printf("∘ %v —  (owner)\n", app)
		} else {
			fmt.Printf("∘ %v —  (collaborator - owned by %s)\n", app, app.Owner.Email)
		}
	}

	return nil
}
