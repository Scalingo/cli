package apps

import (
	"context"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo/v8"
	"github.com/Scalingo/go-scalingo/v8/debug"
	"github.com/Scalingo/go-utils/errors/v2"
)

func Info(ctx context.Context, appName string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Notef(ctx, err, "get Scalingo client")
	}

	app, err := c.AppsShow(ctx, appName)
	if err != nil {
		return errors.Notef(ctx, err, "get the application information from the API")
	}

	stackName, err := getStackName(ctx, c, app.StackID)
	if err != nil {
		debug.Println("Failed to get the stack name from its ID:", err)
		stackName = app.StackID
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.Header([]string{"Settings", "Value"})
	_ = t.Append([]string{"Project", app.ProjectSlug()})
	_ = t.Append([]string{"Force HTTPS", fmt.Sprintf("%v", app.ForceHTTPS)})
	_ = t.Append([]string{"Sticky Session", fmt.Sprintf("%v", app.StickySession)})
	_ = t.Append([]string{"Stack", stackName})
	_ = t.Append([]string{"Status", fmt.Sprintf("%v", app.Status)})
	_ = t.Append([]string{"HDS", fmt.Sprintf("%v", app.HDSResource)})
	_ = t.Render()

	return nil
}

func getStackName(ctx context.Context, c *scalingo.Client, stackID string) (string, error) {
	stacks, err := c.StacksList(ctx)
	if err != nil {
		return "", err
	}

	for _, stack := range stacks {
		if stack.ID == stackID {
			return stack.Name, nil
		}
	}
	return "", errors.Newf(ctx, "unknown stack '%v'", stackID)
}
