package domains

import (
	"context"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v4"
)

func Remove(ctx context.Context, app string, domain string) error {
	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	d, err := findDomain(ctx, client, app, domain)
	if err != nil {
		return errgo.Mask(err)
	}

	err = client.DomainsRemove(ctx, app, d.ID)
	if err != nil {
		return errgo.Mask(err)
	}

	io.Status("The domain", d.Name, "has been deleted")
	return nil
}

func findDomain(ctx context.Context, client *scalingo.Client, app string, domain string) (scalingo.Domain, error) {
	domains, err := client.DomainsList(ctx, app)
	if err != nil {
		return scalingo.Domain{}, errgo.Mask(err)
	}

	for _, d := range domains {
		if d.Name == domain {
			return d, nil
		}
	}
	return scalingo.Domain{}, errgo.New("There is no such domain, please ensure you've added it correctly.\nhttps://my.scalingo.com/apps/" + app + "/domains")
}
