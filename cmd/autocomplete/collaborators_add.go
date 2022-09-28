package autocomplete

import (
	"fmt"
	"sync"

	"github.com/urfave/cli/v2"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo/v6"
	"github.com/Scalingo/go-scalingo/v6/debug"
)

func CollaboratorsAddAutoComplete(c *cli.Context) error {
	var err error
	appName := CurrentAppCompletion(c)
	if appName == "" {
		return nil
	}

	apps, err := appsList(c.Context)
	if err != nil {
		debug.Println("fail to get apps list:", err)
		return nil
	}

	client, err := config.ScalingoClient(c.Context)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	currentAppCollaborators, err := client.CollaboratorsList(c.Context, appName)
	if err != nil {
		return nil
	}

	var apiError error = nil
	ch := make(chan string)
	var wg sync.WaitGroup
	wg.Add(len(apps))
	for _, app := range apps {
		go func(app *scalingo.App) {
			defer wg.Done()
			appCollaborators, erro := client.CollaboratorsList(c.Context, app.Name)
			if erro != nil {
				config.C.Logger.Println(erro.Error())
				apiError = erro
				return
			}
			for _, col := range appCollaborators {
				ch <- col.Email
			}
		}(app)
	}

	setEmails := make(map[string]bool)
	go func() {
		for content := range ch {
			setEmails[content] = true
		}
	}()
	wg.Wait()
	close(ch)

	if apiError != nil {
		return nil
	}

	for email := range setEmails {
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
