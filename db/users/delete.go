package users

import (
	"context"
	"strings"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v6"
	"github.com/Scalingo/go-utils/errors/v2"
)

func DeleteUser(ctx context.Context, app, addonUUID, username string) error {
	isSupported, err := doesDatabaseHandleUserManagement(ctx, app, addonUUID)
	if err != nil {
		return errors.Wrap(ctx, err, "get user management information")
	}

	if !isSupported {
		io.Error(ErrDatabaseNotSupportUserManagement)
		return nil
	}

	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	databaseUsers, err := c.DatabaseListUsers(ctx, app, addonUUID)
	if err != nil {
		return errors.Wrap(ctx, err, "list the database's users")
	}

	var givenUser *scalingo.DatabaseUser
	for _, user := range databaseUsers {
		if strings.EqualFold(user.Name, username) {
			givenUser = &user
			break
		}
	}
	if givenUser == nil {
		io.Error("Confirmation failed: usernames do not match.")
		return nil
	}

	if givenUser.Protected {
		io.Error("Protected user is managed by the platform and cannot be deleted.")
		return nil
	}

	err = c.DatabaseDeleteUser(ctx, app, addonUUID, username)
	if err != nil {
		return errors.Wrapf(ctx, err, "delete given user from database %v", username)
	}
	return nil
}
