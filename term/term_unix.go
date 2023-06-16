//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd
// +build darwin dragonfly freebsd linux netbsd openbsd

package term

import (
	"os"
	"os/exec"
	"strconv"
	"strings"

	"gopkg.in/errgo.v1"
)

// IsTerminal returns true if f is a terminal.
func IsTerminal(f *os.File) bool {
	cmd := exec.Command("test", "-t", "0")
	cmd.Stdin = f
	return cmd.Run() == nil
}

func IsATTY(f *os.File) bool {
	stat, err := f.Stat()
	if err != nil {
		return false
	}

	return (stat.Mode() & os.ModeCharDevice) != 0
}

func MakeRaw(f *os.File) error {
	return makeRaw(f)
}

func Restore(f *os.File) error {
	return restore(f)
}

func Cols() (int, error) {
	cols, err := tput("cols")
	if err != nil {
		return 0, errgo.Mask(err, errgo.Any)
	}
	return strconv.Atoi(cols)
}

func Lines() (int, error) {
	cols, err := tput("lines")
	if err != nil {
		return 0, errgo.Mask(err, errgo.Any)
	}
	return strconv.Atoi(cols)
}

func tput(what string) (string, error) {
	c := exec.Command("tput", what)
	c.Stderr = os.Stderr
	out, err := c.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
