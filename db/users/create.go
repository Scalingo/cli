package users

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/term"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v6"
	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/gopassword"
)

func CreateUser(ctx context.Context, app, addonUUID, username string, readonly bool) error {
	isSupported, err := doesDatabaseHandleUserManagement(ctx, app, addonUUID)
	if err != nil {
		return errors.Wrap(ctx, err, "get user management information")
	}

	if !isSupported {
		io.Error(ErrDatabaseNotSupportUserManagement)
		return nil
	}

	if usernameValidation, ok := isUsernameValid(username); !ok {
		io.Error(usernameValidation)
		return nil
	}

	password, confirmedPassword, err := askForPassword(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "ask for password")
	}

	passwordValidation, ok := isPasswordValid(password, confirmedPassword)
	if !ok && password != "" {
		io.Error(passwordValidation)
		return nil
	}

	isPasswordGenerated := false
	if password == "" {
		isPasswordGenerated = true
		password = gopassword.Generate(64)
		confirmedPassword = password
	}

	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	user := scalingo.DatabaseCreateUserParam{
		DatabaseID:           addonUUID,
		Name:                 username,
		Password:             password,
		PasswordConfirmation: confirmedPassword,
		ReadOnly:             readonly,
	}
	databaseUsers, err := c.DatabaseCreateUser(ctx, app, addonUUID, user)
	if err != nil {
		return errors.Wrap(ctx, err, "create the given database user")
	}

	if isPasswordGenerated {
		fmt.Printf("User \"%s\" created with password \"%s\".\n", databaseUsers.Name, password)
		return nil
	}

	fmt.Printf("User \"%s\" created.\n", databaseUsers.Name)
	return nil
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
	if password != confirmedPassword {
		return "Password confirmation doesn't match", false
	}
	if len(password) < 8 || len(password) > 64 {
		return "Password must contain between 8 and 64 characters", false
	}
	return "", true
}

func isUsernameValid(username string) (string, bool) {
	if len(username) < 6 || len(username) > 32 {
		return "name must contain between 6 and 32 characters", false
	}
	return "", true
}
