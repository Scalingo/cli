package session

import (
	"fmt"

	"github.com/Scalingo/cli/api"
	"github.com/Scalingo/cli/term"
	"gopkg.in/errgo.v1"
)

func SignUp() error {
	fmt.Print("Email: ")

	var email string
	_, err := fmt.Scanln(&email)
	if err != nil {
		return errgo.Mask(err)
	}

	password, err := term.Password("Password: ")
	if err != nil {
		return errgo.Mask(err)
	}

	password_confirmation, err := term.Password("Password validation: ")
	if err != nil {
		return errgo.Mask(err)
	}

	if password != password_confirmation {
		return errgo.New("passwords don't match")
	}

	err = api.SignUp(email, password)
	if err != nil {
		return errgo.Mask(err)
	}

	fmt.Println("A confirmation email has been sent to", email)
	return nil
}
