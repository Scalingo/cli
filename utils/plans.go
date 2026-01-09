package utils

import (
	"context"

	"github.com/Scalingo/go-utils/errors/v2"

	"github.com/Scalingo/go-scalingo/v9"
)

func FindPlan(ctx context.Context, c *scalingo.Client, addon, plan string) (string, error) {
	plans, err := c.AddonProviderPlansList(ctx, addon, scalingo.AddonProviderPlansListOpts{})
	if err != nil {
		return "", errors.Wrapf(ctx, err, "list addon %s plans", addon)
	}
	for _, p := range plans {
		if plan == p.Name {
			return p.ID, nil
		}
	}
	return "", errors.Newf(ctx, "no plan %s found for addon %s", plan, addon)
}
