package utils

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-scalingo/v11/http"
	"github.com/Scalingo/go-utils/errors/v3"
)

const (
	Success = "✔"
	Error   = "✘"
)

// Ask the user whether or not he wants to break his free trial. If not, return without doing
// anything. If yes, call the given callback function.
func AskAndStopFreeTrial(ctx context.Context, c *scalingo.Client, callback func() error) error {
	validate, err := askUserValidation(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to ask for user to validate to break out of free trial")
	}
	if !validate {
		fmt.Println("Do not break free trial.")
		return nil
	}
	err = c.UserStopFreeTrial(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to stop user free trial")
	}
	return callback()
}

// Return true if the given error is because of a Payment Required error and the free trial is
// exceeded.
func IsPaymentRequiredAndFreeTrialExceededError(err error) bool {
	var reqestFailedError *http.RequestFailedError
	if !errors.As(err, &reqestFailedError) || reqestFailedError.Code != 402 {
		return false
	}
	paymentRequiredErr, ok := reqestFailedError.APIError.(http.PaymentRequiredError)
	if !ok || !strings.HasSuffix(paymentRequiredErr.Name, "free-trial-exceeded") {
		return false
	}
	return true
}

func askUserValidation(ctx context.Context) (bool, error) {
	fmt.Println("You are still in your free trial. If you continue, your free trial will end and you will be billed for your usage of the platform. Do you agree? [Y/n]")
	in, err := readCharFromStdin(ctx)
	if err != nil {
		return false, errors.Wrap(ctx, err, "fail to read user confirmation")
	}
	if in != "" && !strings.EqualFold(in, "Y") {
		return false, nil
	}
	return true, nil
}

// Read a single character on stdin. The string is trimmed of white space.
// If the string is then empty, its value is "Y"
func readCharFromStdin(ctx context.Context) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", errors.Wrap(ctx, err, "read confirmation from stdin")
	}
	input = strings.TrimSpace(input)
	if input == "" {
		return "", nil
	}
	return string(input[0]), nil
}
