package term

import (
	"fmt"

	"github.com/howeyc/gopass"

	"gopkg.in/errgo.v1"
)

func Password(prompt string) (string, error) {
	fmt.Printf(prompt)
	bytePassword, err := gopass.GetPasswd()
	if err != nil {
		return "", errgo.Mask(err)
	} else {
		return string(bytePassword), nil
	}
}
