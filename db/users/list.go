package users

import (
	"context"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-utils/errors/v2"
)

func List(ctx context.Context, app, addonUUID string) error {
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

	header := []string{"Username", "Read-Only", "Protected"}
	if len(databaseUsers) > 0 && databaseUsers[0].DbmsAttributes != nil {
		header = append(header, "Password Encryption")
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.Header(header)

	for _, user := range databaseUsers {
		line := []string{
			user.Name,
			fmt.Sprintf("%v", user.ReadOnly),
			fmt.Sprintf("%v", user.Protected),
		}
		if user.DbmsAttributes != nil {
			line = append(line, user.DbmsAttributes.PasswordEncryption)
		}
		t.Append(line)
	}
	t.Render()

	return nil
}
