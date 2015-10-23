package session

import (
	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/config"
)

func Login() error {
	_, err := config.Authenticator.LoadAuth()
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	return nil
}
