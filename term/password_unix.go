// +build darwin dragonfly freebsd linux netbsd openbsd

package term

import (
	"fmt"

	"gopkg.in/errgo.v1"

	"github.com/howeyc/gopass"
)

func Password(prompt string) (string, error) {
	fmt.Printf(prompt)
	res, err := gopass.GetPasswd()
	if err != nil {
		return "", errgo.Mask(err)
	} else {
		return string(res), nil
	}
}
