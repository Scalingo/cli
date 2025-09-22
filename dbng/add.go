package dbng

import (
	"context"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v8"
	"github.com/Scalingo/go-utils/errors/v2"
)

const dbngTypePostgresql = "dedicated-pg"

var aliases = map[string]string{
	"psql":             dbngTypePostgresql,
	"pgsql":            dbngTypePostgresql,
	"postgres":         dbngTypePostgresql,
	"postgresql":       dbngTypePostgresql,
	dbngTypePostgresql: dbngTypePostgresql,
}

func Add(ctx context.Context, params scalingo.DatabaseCreateParams) error {

	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	addon, ok := aliases[params.AddonProviderID]
	if !ok {
		return errors.New(ctx, "invalid database type")
	}

	planID, err := utils.CheckPlanExist(ctx, c, addon, params.PlanID)
	if err != nil {
		return errors.Wrap(ctx, err, "invalid database plan")
	}

	// Adjust params fields
	params.AddonProviderID = addon
	params.PlanID = planID

	db, err := c.Preview().DatabaseCreate(ctx, params)
	if err != nil {
		return errors.Wrap(ctx, err, "add database")
	}

	io.Status(db.App.Name, "has been created")

	return nil
}
