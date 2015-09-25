package autocomplete

import (
	"fmt"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/go-scalingo"
	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
)

func ScaleAutoComplete(c *cli.Context) error {
	appName := CurrentAppCompletion(c)
	if appName == "" {
		return nil
	}

	processes, err := scalingo.AppsPs(appName)
	if err != nil {
		return errgo.Mask(err)
	}
	for _, ct := range processes {
		fmt.Println(fmt.Sprintf("%s:%d:%s", ct.Name, ct.Amount, ct.Size))
	}

	return nil
}
