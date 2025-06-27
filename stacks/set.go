package stacks

import (
	"context"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v8"
)

func Set(ctx context.Context, app string, stack string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	stacks, err := c.StacksList(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to list available stacks")
	}

	var stackToSet scalingo.Stack
	for _, apistack := range stacks {
		if apistack.Name == stack || apistack.ID == stack {
			stackToSet = apistack
			break
		}
	}

	if stackToSet.ID == "" {
		return errgo.Newf("stack '%v' is unknown.", stack)
	}

	_, err = c.AppsSetStack(ctx, app, stackToSet.ID)
	if err != nil {
		return errgo.Notef(err, "fail to set stack %v (%v)", stackToSet.Name, stackToSet.ID)
	}

	io.Statusf("Stack of %v has been set to %v (%v)\n", io.Bold(app), io.Bold(stackToSet.Name), stackToSet.ID)
	io.Infof(io.Gray("Deployment cache of %v has been reset\n"), app)

	return nil
}
