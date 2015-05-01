package apps

import (
	"fmt"

	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/api"
)

func Restart(app string, sync bool, args []string) error {
	params := api.AppsRestartParams{Scope: args}

	res, err := api.AppsRestart(app, &params)
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
