package stacks

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo"
	"gopkg.in/errgo.v1"
)

func Set(app string, stack string) error {
	c := config.ScalingoClient()
	stacks, err := c.StacksList()
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

	_, err = c.AppsSetStack(app, stackToSet.ID)
	if err != nil {
		return errgo.Notef(err, "fail to set stack %v (%v)", stackToSet.Name, stackToSet.ID)
	}

	io.Statusf("Stack of %v has been set to %v (%v)\n", app, stackToSet.Name, stackToSet.ID)

	return nil

}
