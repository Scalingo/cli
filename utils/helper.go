package utils

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/go-scalingo/v7"
	"github.com/Scalingo/go-scalingo/v7/http"
	errors "github.com/Scalingo/go-utils/errors/v2"
)

const (
	Success = "✔"
	Error   = "✘"
)

// Ask the user whether or not he wants to break his free trial. If not, return without doing
// anything. If yes, call the given callback function.
func AskAndStopFreeTrial(ctx context.Context, c *scalingo.Client, callback func() error) error {
	validate, err := askUserValidation()
	if err != nil {
		return errgo.Notef(err, "fail to ask for user to validate to break out of free trial")
	}
	if !validate {
		fmt.Println("Do not break free trial.")
		return nil
	}
	err = c.UserStopFreeTrial(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to stop user free trial")
	}
	return callback()
}

// Return true if the given error is because of a Payment Required error and the free trial is
// exceeded.
func IsPaymentRequiredAndFreeTrialExceededError(err error) bool {
	reqestFailedError, ok := errors.RootCause(err).(*http.RequestFailedError)
	if !ok || reqestFailedError.Code != 402 {
		return false
	}
	paymentRequiredErr, ok := reqestFailedError.APIError.(http.PaymentRequiredError)
	if !ok || !strings.HasSuffix(paymentRequiredErr.Name, "free-trial-exceeded") {
		return false
	}
	return true
}

func askUserValidation() (bool, error) {
	fmt.Println("You are still in your free trial. If you continue, your free trial will end and you will be billed for your usage of the platform. Do you agree? [Y/n]")
	in, err := readCharFromStdin()
	if err != nil {
		return false, errgo.Mask(err, errgo.Any)
	}
	if in != "" && !strings.EqualFold(in, "Y") {
		return false, nil
	}
	return true, nil
}

// Read a single character on stdin. The string is trimmed of white space.
// If the string is then empty, its value is "Y"
func readCharFromStdin() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	input = strings.TrimSpace(input)
	if input == "" {
		return "", nil
	}
	return string(input[0]), nil
}
