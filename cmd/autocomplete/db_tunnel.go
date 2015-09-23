package autocomplete

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/go-scalingo"
)

func DbTunnelAutoComplete(c *cli.Context) error {
	appName := CurrentAppCompletion(c)
	if appName == "" {
		return nil
	}

	lastArg := ""
	if len(os.Args) > 2 {
		lastArg = os.Args[len(os.Args)-2]
	}

	if !strings.HasPrefix(lastArg, "-") {
		variables, err := scalingo.VariablesList(appName)
		if err == nil {

			for _, v := range variables {
				if matched, err := regexp.Match("SCALINGO_.*_URL", []byte(v.Name)); matched && err == nil {
					fmt.Println(v.Name)
				}
			}
		}
	}

	return nil
}
