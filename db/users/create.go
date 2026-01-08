package users

import (
	"context"
	"fmt"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v9"
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

	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}
	l, err := c.DatabaseListUsers(ctx, app, addonUUID)
	if err != nil {
		return errors.Wrap(ctx, err, "list the database's users")
	}

	// Check if the user exists
	var userExists bool
	for _, user := range l {
		if user.Name == username {
			userExists = true
			break
		}
	}
	if userExists {
		return errors.New(ctx, fmt.Sprintf("User \"%s\" already exists", username))
	}

	password, confirmedPassword, err := askForPasswordWithRetry(ctx, 3)
	if err != nil {
		io.Error(err)
		return nil
	}

	isPasswordGenerated := false
	if password == "" {
		isPasswordGenerated = true
		password = gopassword.Generate()
		confirmedPassword = password
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
