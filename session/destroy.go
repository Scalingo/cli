package session

import (
	"github.com/Scalingo/cli/config"
	"gopkg.in/errgo.v1"
)

func DestroyToken() error {
	authenticator := &config.CliAuthenticator{}
	if err := authenticator.RemoveAuth(); err != nil {
		return errgo.Mask(err)
	}
	return nil
}
