package env

import (
	"fmt"

	"gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/config"
)

func Display(app string) error {
	c := config.ScalingoClient()
	vars, err := c.VariablesList(app)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	for _, v := range vars {
		fmt.Printf("%s=%s\n", v.Name, v.Value)
	}
	return nil
}
