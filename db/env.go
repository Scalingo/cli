package db

import (
	"net/url"
	"strings"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func dbURL(appName, envVariableName string, urlSchemes []string) (*url.URL, string, string, error) {
	u, err := dbURLFromAPI(appName, envVariableName, urlSchemes)
	if err != nil {
		return nil, "", "", errgo.Mask(err)
	}

	dbURL, err := url.Parse(u)
	if err != nil {
		return nil, "", "", errgo.Newf("%v is not a valid URL", u)
	}

	user, password, err := extractCredentials(dbURL)
	if err != nil {
		return nil, "", "", errgo.Mask(err)
	}

	return dbURL, user, password, nil
}

func dbURLFromAPI(appName, envVariableName string, urlSchemes []string) (string, error) {
	scalingoClient, err := config.ScalingoClient()
	if err != nil {
		return "", errgo.Notef(err, "fail to get Scalingo client to list the variables")
	}

	variables, err := scalingoClient.VariablesListWithoutAlias(appName)
	if err != nil {
		return "", errgo.Mask(err)
	}
	for _, variable := range variables {
		for _, scheme := range urlSchemes {
			if strings.Contains(variable.Name, envVariableName) && strings.HasPrefix(variable.Value, scheme+"://") {
				return variable.Value, nil
			}
		}
	}

	return "", errgo.Newf("no %v addon detected", strings.ToLower(envVariableName))
}

func extractCredentials(u *url.URL) (string, string, error) {
	if u.User == nil {
		return "", "", errgo.Newf("%v has no information about the instance credentials", u)
	}

	password, ok := u.User.Password()
	if !ok {
		return "", "", errgo.Newf("%v has no information about the instance password", u)
	}

	return u.User.Username(), password, nil
}
