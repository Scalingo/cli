package apps

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"gopkg.in/errgo.v1"
)

func Rename(appName string, newName string) error {
	var validationName string

	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	_, err = c.AppsShow(appName)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	fmt.Printf("/!\\ You're going to rename '%s' to '%s'\nTo confirm type the name of the application: ", appName, newName)
	validationName, err = bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	validationName = strings.Trim(validationName, "\n")

	if validationName != appName {
		return errgo.Newf("'%s' is not '%s', aborting…\n", validationName, appName)
	}

	_, err = c.AppsRename(appName, newName)
	if err != nil {
		return errgo.Notef(err, "fail to rename app")
	}

	io.Status("App " + appName + " has been renamed to " + newName)
	return nil
}
