package dbng

import (
	"context"

	"github.com/Scalingo/cli/addonproviders"
	"github.com/Scalingo/go-utils/errors/v2"
)

func ListPlans(ctx context.Context, technology string) error {
	err := addonproviders.Plans(ctx, technology)
	if err != nil {
		return errors.Wrap(ctx, err, "list the plans of the database")
	}
	return nil
}
