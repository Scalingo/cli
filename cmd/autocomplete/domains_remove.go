package autocomplete

import (
	"fmt"

	"github.com/urfave/cli"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

func DomainsRemoveAutoComplete(c *cli.Context) error {
	appName := CurrentAppCompletion(c)
	if appName == "" {
		return nil
	}

	client, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	domains, err := client.DomainsList(appName)
	if err != nil {
		return errgo.Notef(err, "fail to get domains list")
	}

	for _, domain := range domains {
		fmt.Println(domain.Name)
	}

	return nil
}
