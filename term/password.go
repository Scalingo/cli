package term

import (
	"fmt"
	"os"

	"golang.org/x/term"
	"gopkg.in/errgo.v1"
)

func Password(prompt string) (string, error) {
	fmt.Print(prompt)
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", errgo.Notef(err, "fail to read the password on stdin")
	}

	return string(bytePassword), nil
}
