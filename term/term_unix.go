//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd

package term

import (
	"context"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/Scalingo/go-utils/errors/v3"
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

func Cols(ctx context.Context) (int, error) {
	cols, err := tput(ctx, "cols")
	if err != nil {
		return 0, errors.Wrap(ctx, err, "fail to get terminal columns")
	}
	return strconv.Atoi(cols)
}

func Lines(ctx context.Context) (int, error) {
	cols, err := tput(ctx, "lines")
	if err != nil {
		return 0, errors.Wrap(ctx, err, "fail to get terminal lines")
	}
	return strconv.Atoi(cols)
}

func tput(ctx context.Context, what string) (string, error) {
	c := exec.Command("tput", what)
	c.Stderr = os.Stderr
	out, err := c.Output()
	if err != nil {
		return "", errors.Wrapf(ctx, err, "run tput %s", what)
	}
	return strings.TrimSpace(string(out)), nil
}
