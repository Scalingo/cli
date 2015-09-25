package apps

import (
	"fmt"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/go-scalingo"
	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/io"
)

func Destroy(appName string) error {
	var validationName string
	fmt.Printf("/!\\ You're going to delete %s, this operation is irreversible.\nTo confirm type the name of the application: ", appName)
	fmt.Scan(&validationName)
	if validationName != appName {
		return errgo.Newf("'%s' is not '%s', aborting…\n", validationName, appName)
	}

	res, err := scalingo.AppsDestroy(appName, validationName)
	if err != nil {
		return errgo.Notef(err, "fail to destroy app")
	}
	defer res.Body.Close()

	io.Status("App " + appName + " has been deleted")
	return nil
}
