package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	scalingo "github.com/Scalingo/go-scalingo"

	"gopkg.in/errgo.v1"
)

func AskAndStopFreeTrial(c *scalingo.Client) (bool, error) {
	validate, err := AskUserValidation()
	if err != nil {
		return false, errgo.Mask(err, errgo.Any)
	}
	if !validate {
		fmt.Println("Do not break free trial.")
		return false, nil
	}
	_, err = c.UpdateUser(scalingo.UpdateUserParams{StopFreeTrial: true})
	if err != nil {
		return false, errgo.Mask(err, errgo.Any)
	}
	return true, nil
}

func PaymentRequiredAndFreeTrialExceeded(err error) bool {
	reqestFailedError, ok := errgo.Cause(err).(*scalingo.RequestFailedError)
	if !ok || reqestFailedError.Code != 402 {
		return false
	}
	paymentRequiredErr, ok := reqestFailedError.APIError.(scalingo.PaymentRequiredError)
	if !ok || paymentRequiredErr.Name != "free-trial-exceeded" {
		return false
	}
	return true
}

func AskUserValidation() (bool, error) {
	fmt.Println("You are still in your free trial. If you continue, your free trial will end and you will be billed for your usage of the platform. Do you agree? [Y/n]")
	in, err := readCharFromStdin()
	if err != nil {
		return false, errgo.Mask(err, errgo.Any)
	}
	if in != "" && strings.ToUpper(in) != "Y" {
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
	if len(input) == 0 {
		return "", nil
	}
	return string(input[0]), nil
}
