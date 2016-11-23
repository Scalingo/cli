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

func Destroy(appName string) error {
	var validationName string

	c := config.ScalingoClient()
	_, err := c.AppsShow(appName)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	fmt.Printf("/!\\ You're going to delete %s, this operation is irreversible.\nTo confirm type the name of the application: ", appName)
	validationName, err = bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	validationName = strings.Trim(validationName, "\n")

	if validationName != appName {
		return errgo.Newf("'%s' is not '%s', aborting…\n", validationName, appName)
	}

	err = c.AppsDestroy(appName, validationName)
	if err != nil {
		return errgo.Notef(err, "fail to destroy app")
	}

	io.Status("App " + appName + " has been deleted")
	return nil
}
