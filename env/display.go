package env

import (
	"fmt"

	"github.com/Scalingo/cli/config"
	"gopkg.in/errgo.v1"
)

func Display(app string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	vars, err := c.VariablesList(app)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	for _, v := range vars {
		fmt.Printf("%s=%s\n", v.Name, v.Value)
	}
	return nil
}
