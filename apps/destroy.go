package apps

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Scalingo/go-utils/errors/v2"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
)

func Destroy(ctx context.Context, appName string, force bool) error {
	var validationName string

	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}

	_, err = c.AppsShow(ctx, appName)
	if err != nil {
		return errors.Wrap(ctx, err, "operation failed")
	}

	if !force {
		fmt.Printf("/!\\ You're going to delete %s, this operation is irreversible.\nTo confirm type the name of the application: ", appName)
		validationName, err = bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			return errors.Wrap(ctx, err, "operation failed")
		}
		validationName = strings.Trim(validationName, "\n")

		if validationName != appName {
			return errors.Newf(ctx, "'%s' is not '%s', aborting…\n", validationName, appName)
		}
	} else {
		validationName = appName
	}

	err = c.AppsDestroy(ctx, appName, validationName)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to destroy app")
	}

	io.Status("App " + appName + " has been deleted")
	return nil
}
