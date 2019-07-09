package apps

import (
	"fmt"
	"os"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo"
	"github.com/Scalingo/go-scalingo/debug"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"gopkg.in/errgo.v1"
)

func Info(appName string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	app, err := c.AppsShow(appName)
	if err != nil {
		return errgo.Notef(err, "fail to get the application information")
	}

	stackName, err := getStackName(c, app.StackID)
	if err != nil {
		debug.Println("Failed to get the stack name from its ID:", err)
		stackName = app.StackID
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"Settings", "Value"})
	t.Append([]string{"Force HTTPS", fmt.Sprintf("%v", app.ForceHTTPS)})
	t.Append([]string{"Sticky Session", fmt.Sprintf("%v", app.StickySession)})
	t.Append([]string{"Stack", stackName})
	t.Append([]string{"Status", fmt.Sprintf("%v", app.Status)})
	t.Render()

	return nil
}

func getStackName(c *scalingo.Client, stackID string) (string, error) {
	stacks, err := c.StacksList()
	if err != nil {
		return "", err
	}

	for _, stack := range stacks {
		if stack.ID == stackID {
			return stack.Name, nil
		}
	}
	return "", errors.New("unknown stack")
}
