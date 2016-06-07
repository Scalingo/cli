package apps

import (
	"fmt"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/go-scalingo"
	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/config"
)

func Restart(app string, sync bool, args []string) error {
	params := scalingo.AppsRestartParams{Scope: args}

	c := config.ScalingoClient()
	res, err := c.AppsRestart(app, &params)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	res.Body.Close()

	if !sync {
		fmt.Println("Your application is being restarted.")
		return nil
	}

	err = handleOperation(app, res)
	if err != nil {
		return errgo.Mask(err)
	}

	fmt.Println("Your application has been restarted.")
	return nil
}
