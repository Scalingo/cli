package apps

import (
	"github.com/Scalingo/cli/api"
	"gopkg.in/errgo.v1"
)

func Restart(app string, args []string) error {
	scope := api.AppRestartScope(args)

	res, err := api.AppsRestart(app, scope)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	res.Body.Close()

	return nil
}
