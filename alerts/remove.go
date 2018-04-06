package alerts

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"gopkg.in/errgo.v1"
)

func Remove(app, id string) error {
	err := config.ScalingoClient().AlertRemove(app, id)
	if err != nil {
		return errgo.Mask(err)
	}

	io.Status("The alert has been deleted")
	return nil
}
