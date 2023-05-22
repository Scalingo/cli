package env

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"os"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo/v6"
	"github.com/joho/godotenv"
)

var (
	setInvalidSyntaxError = errors.New("format is invalid, accepted: VAR=VAL")

	nameFormat           = regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)
	errInvalidNameFormat = fmt.Errorf("name can only be composed with alphanumerical characters, hyphens and underscores")
)

func Add(ctx context.Context, app string, params []string, filePath string) error {
	var variables scalingo.Variables

	err := readFromFile(&variables, filePath)
	if err != nil {
		return err
	}

	err = readFromCmdLine(&variables, params)
	if err != nil {
		return err
	}

	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	_, _, err = c.VariableMultipleSet(ctx, app, variables)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	for _, variable := range variables {
		fmt.Printf("%s has been set to '%s'.\n", variable.Name, variable.Value)
	}
	fmt.Println("\nRestart your containers to apply these environment changes on your application:")
	fmt.Printf("scalingo --app %s restart\n", app)

	return nil
}

func Delete(ctx context.Context, app string, varNames []string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	vars, err := c.VariablesList(ctx, app)

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
		err := c.VariableUnset(ctx, app, v.ID)
		if err != nil {
			return errgo.Mask(err, errgo.Any)
		}
		fmt.Printf("%s has been unset.\n", v.Name)
	}
	fmt.Println("\nRestart your containers to apply these environment changes on your application:")
	fmt.Printf("scalingo --app %s restart\n", app)
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

	if !nameFormat.MatchString(name) {
		return errInvalidNameFormat
	}

	return nil
}

func parseVariable(param string) (string, string) {
	editSplit := strings.SplitN(param, "=", 2)
	return editSplit[0], editSplit[1]
}

func readFromCmdLine(variables *scalingo.Variables, params []string) error {
	for _, param := range params {
		if err := isEnvEditValid(param); err != nil {
			return errgo.Newf("'%s' is invalid: %s", param, err)
		}

		name, value := parseVariable(param)
		*variables = append(*variables, &scalingo.Variable{
			Name:  name,
			Value: value,
		})
	}
	return nil
}

func readFromFile(variables *scalingo.Variables, filePath string) error {
	if len(filePath) > 0 {
		var env map[string]string
		var err error
		if filePath == "-" {
			env, err = godotenv.Parse(os.Stdin)
			if err != nil {
				return errgo.Newf("Error while reading from stdin: %s", err)
			}
		} else {
			env, err = godotenv.Read(filePath)
			if err != nil {
				return errgo.Newf("File is invalid: %s", err)
			}
		}
		for name, value := range env {
			*variables = append(*variables, &scalingo.Variable{
				Name:  name,
				Value: value,
			})
		}
	}
	return nil
}
