package env

import (
	"fmt"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/go-scalingo"
	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
)

func Display(app string) error {
	vars, err := scalingo.VariablesList(app)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	for _, v := range vars {
		fmt.Printf("%s=%s\n", v.Name, v.Value)
	}
	return nil
}
