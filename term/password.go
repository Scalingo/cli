package term

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/term"

	"github.com/Scalingo/go-utils/errors/v3"
)

func Password(prompt string) (string, error) {
	fmt.Print(prompt)
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", errors.Wrapf(context.Background(), err, "fail to read the password on stdin")
	}

	return string(bytePassword), nil
}
