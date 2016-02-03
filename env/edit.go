package env

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/go-scalingo"
	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
)

var (
	setInvalidSyntaxError = errors.New("format is invalid, accepted: VAR=VAL")
	valueTooLongError     = fmt.Errorf("value is too long (max %d)", scalingo.EnvValueMaxLength)
	nameTooLongError      = fmt.Errorf("name is too long (max %d)", scalingo.EnvNameMaxLength)

	nameFormat             = regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)
	invalidNameFormatError = fmt.Errorf("name can only be composed with alphanumerical characters, hyphens and underscores")
)

func Add(app string, params []string) error {
	var variables scalingo.Variables
	for _, param := range params {
		if err := isEnvEditValid(param); err != nil {
			return errgo.Newf("'%s' is invalid: %s", param, err)
		}

		name, value := parseVariable(param)
		variables = append(variables, &scalingo.Variable{
			Name:  name,
			Value: value,
		})
	}

	_, _, err := scalingo.VariableMultipleSet(app, variables)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	for _, variable := range variables {
		fmt.Printf("%s has been set to '%s'.\n", variable.Name, variable.Value)
	}

	return nil
}

func Delete(app string, varNames []string) error {
	vars, err := scalingo.VariablesList(app)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	var varsToUnset scalingo.Variables

	for _, varName := range varNames {
		v, ok := vars.Contains(varName)
		if !ok {
			return errgo.Newf("%s variable does not exist", varName)
		}
		varsToUnset = append(varsToUnset, v)
	}

	for _, v := range varsToUnset {
		err := scalingo.VariableUnset(app, v.ID)
		if err != nil {
			return errgo.Mask(err, errgo.Any)
		}
		fmt.Printf("%s has been unset.\n", v.Name)
	}
	return nil
}

func isEnvEditValid(edit string) error {
	if !strings.Contains(edit, "=") {
		return setInvalidSyntaxError
	}
	name, value := parseVariable(edit)

	if name == "" || value == "" {
		return setInvalidSyntaxError
	}

	if len(name) > scalingo.EnvNameMaxLength {
		return nameTooLongError
	}

	if len(value) > scalingo.EnvValueMaxLength {
		return valueTooLongError
	}

	if !nameFormat.MatchString(name) {
		return invalidNameFormatError
	}

	return nil
}

func parseVariable(param string) (string, string) {
	editSplit := strings.SplitN(param, "=", 2)
	return editSplit[0], editSplit[1]
}
