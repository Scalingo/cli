// +build darwin dragonfly freebsd linux netbsd openbsd

package term

import (
	"github.com/Scalingo/cli/Godeps/_workspace/src/code.google.com/p/gopass"
)

func Password(prompt string) (string, error) {
	return gopass.GetPass(prompt)
}
