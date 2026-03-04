package domains

import (
	"context"

	"github.com/Scalingo/go-utils/errors/v3"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
)

func SetCanonical(ctx context.Context, app, domain string) error {
	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}

	d, err := findDomain(ctx, client, app, domain)
	if err != nil {
		return errors.Wrapf(ctx, err, "find domain %s on app %s", domain, app)
	}

	_, err = client.DomainSetCanonical(ctx, app, d.ID)
	if err != nil {
		return errors.Wrapf(ctx, err, "set canonical domain %s on app %s", domain, app)
	}

	io.Statusf("Canonical domain set to %s\n", domain)
	return nil
}

func UnsetCanonical(ctx context.Context, app string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}

	_, err = c.DomainUnsetCanonical(ctx, app)
	if err != nil {
		return errors.Wrapf(ctx, err, "unset canonical domain on app %s", app)
	}

	io.Status("Canonical domain disabled")
	return nil
}
