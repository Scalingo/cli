package session

import (
	"fmt"
	"os"

	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo"
)

func Login() error {
	user, err := scalingo.AuthFromConfig()
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	if user == nil {
		fmt.Fprintln(os.Stderr, "You need to be authenticated to use Scalingo client.\nNo account ? → https://my.scalingo.com/users/signup")
		user, err = config.Auth()
		if err != nil {
			return errgo.Mask(err, errgo.Any)
		}
	} else {
		io.Status("You are already identified as", user.Username, "<"+user.Email+">")
	}
	return nil
}
