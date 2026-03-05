package env

import (
	"context"
	"fmt"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-utils/errors/v3"
)

func Display(ctx context.Context, app string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}
	vars, err := c.VariablesList(ctx, app)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to list the environment variables")
	}

	for _, v := range vars {
		fmt.Printf("%s=%s\n", v.Name, v.Value)
	}
	return nil
}

func Get(ctx context.Context, appName, variableName string) (string, error) {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return "", errors.Wrapf(ctx, err, "fail to get Scalingo client to get an environment variable")
	}
	vars, err := c.VariablesListWithoutAlias(ctx, appName)
	if err != nil {
		return "", errors.Wrapf(ctx, err, "fail to list the environment variables")
	}

	for _, v := range vars {
		if v.Name == variableName {
			return v.Value, nil
		}
	}
	return "", errors.New(ctx, "environment variable not found")
}
