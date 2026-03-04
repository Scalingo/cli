package autocomplete

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/go-utils/errors/v2"

	"github.com/Scalingo/cli/config"
)

func DomainsRemoveAutoComplete(ctx context.Context, c *cli.Command) error {
	appName := CurrentAppCompletion(c)
	if appName == "" {
		return nil
	}

	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}
	domains, err := client.DomainsList(ctx, appName)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get domains list")
	}

	for _, domain := range domains {
		fmt.Println(domain.Name)
	}

	return nil
}
