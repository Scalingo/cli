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

	var msg string
	if params.Disabled != nil {
		if *params.Disabled {
			msg = "Alert disabled"
		} else {
			msg = "Alert enabled"
		}
	} else {
		msg = "Alert updated"
	}
	io.Status(msg)
	return nil
}
