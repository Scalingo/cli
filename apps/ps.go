package apps

import (
	"fmt"

	"github.com/Scalingo/cli/api"
	"github.com/Scalingo/cli/io"
	"gopkg.in/errgo.v1"
)

func Ps(app string) error {
	processes, err := api.AppsPs(app)
	if err != nil {
		return errgo.Mask(err)
	}

	fmt.Printf("Application processes:\n")
	for _, ct := range processes {
		fmt.Println(io.Indent(fmt.Sprintf("%s: %d", ct.Name, ct.Amount), 2))
	}
	return nil
}
