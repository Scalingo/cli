package apps

import (
	"fmt"

	"github.com/Scalingo/cli/api"
	"gopkg.in/errgo.v1"
)

func Restart(app string, sync bool, args []string) error {
	params := api.AppsRestartParams{args}

	res, err := api.AppsRestart(app, &params)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	res.Body.Close()

	if !sync {
		fmt.Println("You application is being restarted.")
		return nil
	}

	err = handleOperation(app, res)
	if err != nil {
		return errgo.Mask(err)
	}

	fmt.Println("Your application has been restarted.")
	return nil
}
