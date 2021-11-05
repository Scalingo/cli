package term

import (
	"fmt"
	"syscall"

	"golang.org/x/term"
	"gopkg.in/errgo.v1"
)

func Password(prompt string) (string, error) {
	fmt.Print(prompt)
	bytePassword, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return "", errgo.Notef(err, "fail to read the password on stdin")
	} else {
		return string(bytePassword), nil
	}
}
