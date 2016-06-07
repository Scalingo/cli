// +build darwin dragonfly freebsd linux netbsd openbsd

package term

import (
	"code.google.com/p/gopass"
)

func Password(prompt string) (string, error) {
	return gopass.GetPass(prompt)
}
