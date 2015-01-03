package session

import (
	"github.com/Scalingo/cli/config"
	"gopkg.in/errgo.v1"
)

func DestroyToken() error {
	if err := config.RemoveAuth(); err != nil {
		return errgo.Mask(err)
	}
	return nil
}
