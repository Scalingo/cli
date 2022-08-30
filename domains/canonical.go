package domains

import (
	"context"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
)

func SetCanonical(ctx context.Context, app, domain string) error {
	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	d, err := findDomain(ctx, client, app, domain)
	if err != nil {
		return errgo.Mask(err)
	}

	_, err = client.DomainSetCanonical(ctx, app, d.ID)
	if err != nil {
		return errgo.Mask(err)
	}

	io.Statusf("Canonical domain set to %s\n", domain)
	return nil
}

func UnsetCanonical(ctx context.Context, app string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	_, err = c.DomainUnsetCanonical(ctx, app)
	if err != nil {
		return errgo.Mask(err)
	}

	io.Status("Canonical domain disabled")
	return nil
}
