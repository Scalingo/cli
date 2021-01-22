package apps

import (
	"fmt"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo/v4"
	"gopkg.in/errgo.v1"
)

func Restart(app string, sync bool, args []string) error {
	params := scalingo.AppsRestartParams{Scope: args}

	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

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
