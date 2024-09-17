package users

import (
	"context"
	stdErrors "errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-utils/errors/v2"
)

var (
	SupportedAddons                     = []string{"PostgreSQL", "InfluxDB", "MongoDB", "MySQL"}
	ErrDatabaseNotSupportUserManagement = stdErrors.New("Error: Addon does not support user management")
)

func doesDatabaseHandleUserManagement(ctx context.Context, app, addonUUID string) (bool, error) {
	addonsClient, err := config.ScalingoClient(ctx)
	if err != nil {
		return false, errors.Wrap(ctx, err, "get Scalingo client")
	}

	addon, err := addonsClient.AddonShow(ctx, app, addonUUID)
	if err != nil {
		return false, errors.Wrap(ctx, err, "get the addon to check user management support")
	}

	for _, supportedAddon := range SupportedAddons {
		if strings.EqualFold(supportedAddon, addon.AddonProvider.Name) {
			return true, nil
		}
	}

	return false, nil
}

func askForPasswordWithRetry(ctx context.Context, remainingRetries int) (string, string, error) {
	if remainingRetries <= 0 {
		return "", "", errors.New(ctx, "No retries left")
	}

	password, confirmedPassword, err := askForPassword(ctx)
	if err != nil {
		return "", "", errors.Wrap(ctx, err, "ask for password")
	}

	passwordValidation, ok := isPasswordValid(password, confirmedPassword)
	if !ok {
		if remainingRetries == 1 {
			return "", "", errors.Newf(ctx, "%s. Too many retries", passwordValidation)
		}
		io.Error(fmt.Sprintf("%s. Remaining retries: %v", passwordValidation, remainingRetries-1))
		return askForPasswordWithRetry(ctx, remainingRetries-1)
	}

	return password, confirmedPassword, nil
}

func askForPassword(ctx context.Context) (string, string, error) {
	fmt.Printf("Password (Will be generated if left empty): ")

	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", "", errors.Wrap(ctx, err, "read password")
	}

	fmt.Printf("\nPassword Confirmation: ")
	confirmedPassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", "", errors.Wrap(ctx, err, "read password confirmation")
	}
	fmt.Println()

	return string(password), string(confirmedPassword), nil
}

func isPasswordValid(password, confirmedPassword string) (string, bool) {
	if password == "" && confirmedPassword == "" {
		return "", true
	}

	if password != confirmedPassword {
		return "Password confirmation doesn't match", false
	}
	if len(password) < 24 || len(password) > 64 {
		return "Password must contain between 24 and 64 characters", false
	}
	return "", true
}

func isUsernameValid(username string) (string, bool) {
	if len(username) < 6 || len(username) > 32 {
		return "Name must contain between 6 and 32 characters", false
	}
	return "", true
}
