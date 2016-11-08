package term

import (
	"fmt"

	"golang.org/x/crypto/ssh/terminal"
	"syscall"

	"gopkg.in/errgo.v1"
)

func Password(prompt string) (string, error) {
	fmt.Printf(prompt)
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", errgo.Mask(err)
	} else {
		return string(bytePassword), nil
	}
}
