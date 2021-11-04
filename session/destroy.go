package session

import (
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func DestroyToken() error {
	authenticator := &config.CliAuthenticator{}
	if err := authenticator.RemoveAuth(); err != nil {
		return errgo.Mask(err)
	}
	return nil
}
