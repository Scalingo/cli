package dbng

import (
	"context"
	"time"

	"github.com/briandowns/spinner"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v8"
	"github.com/Scalingo/go-utils/errors/v2"
)

func Create(ctx context.Context, params scalingo.DatabaseCreateParams, wait bool) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	addon := params.AddonProviderID

	planID, err := utils.FindPlan(ctx, c, params.AddonProviderID, params.PlanID)
	if err != nil {
		return errors.Wrap(ctx, err, "invalid database plan")
	}

	// Adjust params fields.
	params.AddonProviderID = addon
	params.PlanID = planID

	db, err := c.Preview().DatabaseCreate(ctx, params)
	if err != nil {
		return errors.Wrap(ctx, err, "create database")
	}
	io.Statusf("Your %s database %s ('%s') is being provisioned…\n\n",
		db.Technology,
		db.ID,
		db.Name,
	)

	if wait {
		spinner := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		spinner.Suffix = " Waiting for the database to be ready\n"
		spinner.Start()
		defer spinner.Stop()

		err = waitForRunningDatabase(ctx, c, db.ID)
		if err != nil {
			return errors.Wrap(ctx, err, "wait for running database")
		}

		spinner.Stop()
		io.Status("Your database is up and running.")
	}

	return nil
}

// waitForRunningDatabase is a blocking function waiting for the database to be running.
func waitForRunningDatabase(ctx context.Context, client *scalingo.Client, appID string) error {
	for range time.Tick(10 * time.Second) {
		db, err := client.Preview().DatabaseShow(ctx, appID)
		if errors.Is(err, scalingo.ErrDatabaseNotFound) {
			continue
		} else if err != nil {
			return err
		}
		if db.Database.Status == scalingo.DatabaseStatusRunning {
			break
		}
	}
	return nil
}
