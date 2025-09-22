package utils

import (
	"context"

	"github.com/Scalingo/go-utils/errors/v2"

	"github.com/Scalingo/go-scalingo/v8"
)

func CheckPlanExist(ctx context.Context, c *scalingo.Client, addon, plan string) (string, error) {
	plans, err := c.AddonProviderPlansList(ctx, addon)
	if err != nil {
		return "", errors.Wrapf(ctx, err, "list addon %s plans", addon)
	}
	for _, p := range plans {
		if plan == p.Name {
			return p.ID, nil
		}
	}
	return "", errors.Newf(ctx, "plan %s doesn't exist for addon %s", plan, addon)
}
