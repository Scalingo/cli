package autocomplete

import (
	"fmt"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/codegangsta-cli"
	"github.com/Scalingo/cli/api"
)

func DomainsRemoveAutoComplete(c *cli.Context) error {
	appName := CurrentAppCompletion(c)
	if appName == "" {
		return nil
	}

	domains, err := api.DomainsList(appName)
	if err == nil {

		for _, domain := range domains {
			fmt.Println(domain.Name)
		}
	}

	return nil
}
