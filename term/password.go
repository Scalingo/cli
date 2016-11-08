package term

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"

	"gopkg.in/errgo.v1"
)

func Password(prompt string) (string, error) {
	fmt.Printf(prompt)
	bytePassword, err := terminal.ReadPassword(os.Stdin.Fd())
	if err != nil {
		return "", errgo.Mask(err)
	} else {
		return string(bytePassword), nil
	}
}
