package autocomplete

import (
	"fmt"
	"sync"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/cli/api"
	"github.com/Scalingo/cli/debug"
)

func CollaboratorsAddAutoComplete(c *cli.Context) error {
	var err error
	appName := CurrentAppCompletion(c)
	if appName == "" {
		return nil
	}

	apps, err := appsList()
	if err != nil {
		debug.Println("fail to get apps list:", err)
		return nil
	}

	currentAppCollaborators, err := api.CollaboratorsList(appName)
	if err != nil {
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(len(apps))
	setEmails := make(map[string]bool)
	for _, app := range apps {
		go func(app *api.App, setEmails map[string]bool) {
			appCollaborators, erro := api.CollaboratorsList(app.Name)
			if erro != nil {
				err = erro
				return
			}
			for _, col := range appCollaborators {
				setEmails[col.Email] = true
			}
			wg.Done()
		}(app, setEmails)
	}
	wg.Wait()
	if err != nil {
		return nil
	}

	for email, _ := range setEmails {
		isAlreadyCollaborator := false
		for _, currentAppCol := range currentAppCollaborators {
			if currentAppCol.Email == email {
				isAlreadyCollaborator = true
			}
		}
		if !isAlreadyCollaborator {
			fmt.Println(email)
		}
	}
	return nil
}
