package alerts

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo"
	"gopkg.in/errgo.v1"
)

func Update(app, id string, params scalingo.AlertUpdateParams) error {
	_, err := config.ScalingoClient().AlertUpdate(app, id, params)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	io.Status("Alert updated")
	io.Info("http://")
	return nil
}
