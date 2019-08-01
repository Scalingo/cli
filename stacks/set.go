package stacks

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"gopkg.in/errgo.v1"
)

func Set(app string, stack string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	res, err := c.AppsSetStack(app, stack)
	if err != nil {
		return errgo.Notef(err, "fail to set stack %v", stack)
	}

	io.Statusf("Stack of %v has been set to %v (%v)\n", io.Bold(app), io.Bold(stack), res.StackID)
	io.Infof(io.Gray("Deployment cache of %v has been reseted\n"), app)

	return nil
}
