package env

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/Scalingo/cli/api"
)

var (
	setInvalidSyntaxError = errors.New("format is invalid, accepted: VAR=VAL")
	valueTooLongError     = fmt.Errorf("value is too long (max %d)", api.EnvValueMaxLength)
	nameTooLongError      = fmt.Errorf("name is too long (max %d)", api.EnvNameMaxLength)

	nameFormat             = regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)
	invalidNameFormatError = fmt.Errorf("name can only be composed with alphanumerical characters, hyphens and underscores")
)

func Add(app string, params []string) error {
	for _, param := range params {
		if err := isEnvEditValid(param); err != nil {
			return fmt.Errorf("'%s' is invalid: %s", param, err)
		}
	}

	for _, param := range params {
		name, value := parseVariable(param)
		_, code, err := api.VariableSet(app, name, value)
		if err != nil {
			return err
		}

		if code == 201 {
			fmt.Printf("%s has been set to %s.\n", name, value)
		} else if code == 200 {
			fmt.Printf("%s has been updated to %s.\n", name, value)
		} else {
			return fmt.Errorf("invalid return code %v\n", code)
		}
	}

	return nil
}

func Delete(app string, varNames []string) error {
	vars, err := api.VariablesList(app)
	if err != nil {
		return err
	}

	var varsToUnset api.Variables

	for _, varName := range varNames {
		v, ok := vars.Contains(varName)
		if !ok {
			return fmt.Errorf("%s variable does not exist", varName)
		}
		varsToUnset = append(varsToUnset, v)
	}

	for _, v := range varsToUnset {
		err := api.VariableUnset(app, v.ID)
		if err != nil {
			return err
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

	if len(name) > api.EnvNameMaxLength {
		return nameTooLongError
	}

	if len(value) > api.EnvValueMaxLength {
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
