package apps

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Scalingo/go-utils/errors/v3"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
)

func Rename(ctx context.Context, appName string, newName string) error {
	var validationName string

	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}

	_, err = c.AppsShow(ctx, appName)
	if err != nil {
		return errors.Wrapf(ctx, err, "check that app %s exists", appName)
	}

	fmt.Printf("/!\\ You're going to rename '%s' to '%s'\nTo confirm type the name of the application: ", appName, newName)
	validationName, err = bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return errors.Wrap(ctx, err, "read rename confirmation from stdin")
	}
	validationName = strings.Trim(validationName, "\n")

	if validationName != appName {
		return errors.Newf(ctx, "'%s' is not '%s', aborting…\n", validationName, appName)
	}

	_, err = c.AppsRename(ctx, appName, newName)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to rename app")
	}

	io.Status("App " + appName + " has been renamed to " + newName)
	return nil
}
