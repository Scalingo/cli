package env

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo/v7"
	scalingoerrors "github.com/Scalingo/go-utils/errors/v2"
)

var (
	errSetInvalidSyntax = errors.New("format is invalid, accepted: VAR=VAL")

	nameFormat           = regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)
	errInvalidNameFormat = fmt.Errorf("name can only be composed with alphanumerical characters, hyphens and underscores")
)

func Add(ctx context.Context, app string, params []string, filePath string) error {
	variablesFromFile, err := readFromFile(ctx, filePath)
	if err != nil {
		return scalingoerrors.Notef(ctx, err, "read .env file")
	}

	variables, err := readFromCmdLine(ctx, variablesFromFile, params)
	if err != nil {
		return scalingoerrors.Notef(ctx, err, "read variables from command line")
	}

	scalingoVariables := scalingo.Variables{}
	for name, value := range variables {
		scalingoVariables = append(scalingoVariables, &scalingo.Variable{
			Name:  name,
			Value: value,
		})
	}

	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return scalingoerrors.Notef(ctx, err, "get Scalingo client")
	}
	_, _, err = c.VariableMultipleSet(ctx, app, scalingoVariables)
	if err != nil {
		return scalingoerrors.Notef(ctx, err, "set multiple variables")
	}

	for _, variable := range scalingoVariables {
		fmt.Printf("%s has been set to '%s'.\n", variable.Name, variable.Value)
	}
	fmt.Println("\nRestart your containers to apply these environment changes on your application:")
	fmt.Printf("scalingo --app %s restart\n", app)

	return nil
}

func Delete(ctx context.Context, app string, varNames []string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return scalingoerrors.Notef(ctx, err, "get Scalingo client")
	}
	vars, err := c.VariablesList(ctx, app)
	if err != nil {
		return scalingoerrors.Notef(ctx, err, "list variables before deletion")
	}

	var varsToUnset scalingo.Variables

	for _, varName := range varNames {
		v, ok := vars.Contains(varName)
		if !ok {
			return scalingoerrors.Newf(ctx, "%s variable does not exist", varName)
		}
		varsToUnset = append(varsToUnset, v)
	}

	for _, v := range varsToUnset {
		err := c.VariableUnset(ctx, app, v.ID)
		if err != nil {
			return scalingoerrors.Notef(ctx, err, "unset variable %s", v.Name)
		}
		fmt.Printf("%s has been unset.\n", v.Name)
	}
	fmt.Println("\nRestart your containers to apply these environment changes on your application:")
	fmt.Printf("scalingo --app %s restart\n", app)
	return nil
}

func isEnvEditValid(edit string) error {
	if !strings.Contains(edit, "=") {
		return errSetInvalidSyntax
	}
	name, value := parseVariable(edit)

	if name == "" || value == "" {
		return errSetInvalidSyntax
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

func readFromCmdLine(ctx context.Context, variables map[string]string, params []string) (map[string]string, error) {
	for _, param := range params {
		err := isEnvEditValid(param)
		if err != nil {
			return nil, scalingoerrors.Newf(ctx, "'%s' is invalid: %s", param, err)
		}

		name, value := parseVariable(param)
		variables[name] = value
	}

	return variables, nil
}

func readFromFile(ctx context.Context, filePath string) (map[string]string, error) {
	if filePath == "" {
		return map[string]string{}, nil
	}

	if filePath == "-" {
		variables, err := godotenv.Parse(os.Stdin)
		if err != nil {
			return nil, scalingoerrors.Notef(ctx, err, "parse .env from stdin")
		}
		return variables, nil
	}

	variables, err := godotenv.Read(filePath)
	if err != nil {
		return nil, scalingoerrors.Notef(ctx, err, "invalid .env file")
	}
	return variables, nil
}
