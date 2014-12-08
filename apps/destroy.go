package apps

import (
	"fmt"

	"github.com/Scalingo/cli/api"
	"gopkg.in/errgo.v1"
)

func Destroy(id string) error {
	res, err := api.AppsDestroy(id)
	if err != nil {
		return errgo.Notef(err, "fail to create app")
	}
	defer res.Body.Close()

	fmt.Printf("App %s has been deleted\n", id)
	return nil
}
