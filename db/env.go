package db

import (
	"context"
	"net/url"
	"strings"

	"github.com/Scalingo/go-utils/errors/v2"

	"github.com/Scalingo/cli/config"
)

func dbURL(ctx context.Context, appName, envVariableName string, urlSchemes []string) (*url.URL, string, string, error) {
	u, err := dbURLFromAPI(ctx, appName, envVariableName, urlSchemes)
	if err != nil {
		return nil, "", "", errors.Wrapf(ctx, err, "fail to retrieve %s database URL", strings.ToLower(envVariableName))
	}

	dbURL, err := url.Parse(u)
	if err != nil {
		return nil, "", "", errors.Newf(ctx, "%v is not a valid URL", u)
	}

	user, password, err := extractCredentials(dbURL)
	if err != nil {
		return nil, "", "", errors.Wrapf(ctx, err, "fail to extract credentials from %s database URL", strings.ToLower(envVariableName))
	}

	return dbURL, user, password, nil
}

func dbURLFromAPI(ctx context.Context, appName, envVariableName string, urlSchemes []string) (string, error) {
	scalingoClient, err := config.ScalingoClient(ctx)
	if err != nil {
		return "", errors.Wrapf(ctx, err, "fail to get Scalingo client to list the variables")
	}

	variables, err := scalingoClient.VariablesListWithoutAlias(ctx, appName)
	if err != nil {
		return "", errors.Wrap(ctx, err, "operation failed")
	}
	for _, variable := range variables {
		for _, scheme := range urlSchemes {
			if strings.Contains(variable.Name, envVariableName) && strings.HasPrefix(variable.Value, scheme+"://") {
				return variable.Value, nil
			}
		}
	}

	return "", errors.Newf(ctx, "no %v addon detected", strings.ToLower(envVariableName))
}

func extractCredentials(u *url.URL) (string, string, error) {
	if u.User == nil {
		return "", "", errors.Newf(context.Background(), "%v has no information about the instance credentials", u)
	}

	password, ok := u.User.Password()
	if !ok {
		return "", "", errors.Newf(context.Background(), "%v has no information about the instance password", u)
	}

	return u.User.Username(), password, nil
}
