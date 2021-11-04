package alerts

import (
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v4"
)

func Update(app, id string, params scalingo.AlertUpdateParams) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	_, err = c.AlertUpdate(app, id, params)
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
