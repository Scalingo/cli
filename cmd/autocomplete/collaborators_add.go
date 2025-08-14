package autocomplete

import (
	"context"
	"fmt"
	"sync"

	"github.com/urfave/cli/v3"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo/v8"
	"github.com/Scalingo/go-scalingo/v8/debug"
)

func CollaboratorsAddAutoComplete(ctx context.Context, c *cli.Command) error {
	var err error
	appName := CurrentAppCompletion(c)
	if appName == "" {
		return nil
	}

	apps, err := appsList(ctx)
	if err != nil {
		debug.Println("fail to get apps list:", err)
		return nil
	}

	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	currentAppCollaborators, err := client.CollaboratorsList(ctx, appName)
	if err != nil {
		return nil
	}

	var apiError error
	ch := make(chan string)
	var wg sync.WaitGroup
	wg.Add(len(apps))
	for _, app := range apps {
		go func(app *scalingo.App) {
			defer wg.Done()
			appCollaborators, erro := client.CollaboratorsList(ctx, app.Name)
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
