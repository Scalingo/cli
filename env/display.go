package env

import (
	"errors"
	"fmt"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func Display(app string) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	vars, err := c.VariablesList(app)
	if err != nil {
		return errgo.Notef(err, "fail to list the environment variables")
	}

	for _, v := range vars {
		fmt.Printf("%s=%s\n", v.Name, v.Value)
	}
	return nil
}

func Get(appName, variableName string) (string, error) {
	c, err := config.ScalingoClient()
	if err != nil {
		return "", errgo.Notef(err, "fail to get Scalingo client to get an environment variable")
	}
	vars, err := c.VariablesListWithoutAlias(appName)
	if err != nil {
		return "", errgo.Notef(err, "fail to list the environment variables")
	}

	for _, v := range vars {
		if v.Name == variableName {
			return v.Value, nil
		}
	}
	return "", errors.New("environment variable not found")
}
