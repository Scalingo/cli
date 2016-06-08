package session

import (
	"gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/config"
)

func DestroyToken() error {
	if err := config.Authenticator.RemoveAuth(); err != nil {
		return errgo.Mask(err)
	}
	return nil
}
