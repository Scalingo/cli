package session

import (
	"fmt"

	"github.com/Scalingo/cli/api"
	"github.com/Scalingo/cli/io"
	"gopkg.in/errgo.v1"
)

func Login() error {
	user, err := api.AuthFromConfig()
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	if user == nil {
		fmt.Println("You need to be authenticated to use Scalingo client.\nNo account ? → https://my.scalingo.com/users/signup")
		user, err = api.Auth()
		if err != nil {
			return errgo.Mask(err, errgo.Any)
		}
	} else {
		io.Status("You are already identified as", user.Username, "<"+user.Email+">")
	}
	return nil
}
