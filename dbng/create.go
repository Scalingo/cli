package dbng

import (
	"context"
	"time"

	"github.com/briandowns/spinner"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-utils/errors/v3"
	"github.com/Scalingo/go-utils/retry"
)

const (
	databaseStatusPollInterval = 30 * time.Second
	databaseStatusMaxWait      = 7 * 24 * time.Hour
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

	retryer := retry.New(
		retry.WithWaitDuration(databaseStatusPollInterval),
		retry.WithMaxDuration(databaseStatusMaxWait),
		retry.WithoutMaxAttempts(),
	)

	if wait {
		spinner := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		spinner.Suffix = " Waiting for the database to be ready\n"
		spinner.Start()
		defer spinner.Stop()

		err := waitForRunningDatabase(ctx, c.Preview(), retryer, db.ID)
		if err != nil {
			return errors.Wrap(ctx, err, "wait for running database")
		}

		spinner.Stop()
		io.Status("Your database is up and running.")
	}

	return nil
}

// waitForRunningDatabase waits until the database is running, the context is
// canceled, or the provisioning timeout is reached.
func waitForRunningDatabase(ctx context.Context, previewClient *scalingo.PreviewClient, retryer retry.Retry, appID string) error {
	err := retryer.Do(ctx, func(ctx context.Context) error {
		db, err := previewClient.DatabaseShow(ctx, appID)
		if errors.Is(err, scalingo.ErrDatabaseNotFound) {
			return errors.Wrapf(ctx, err, "get database %s status", appID)
		}
		if err != nil {
			return retry.NewRetryCancelError(errors.Wrapf(ctx, err, "get database %s status", appID))
		}

		switch db.Database.Status {
		case scalingo.DatabaseStatusRunning:
			return nil
		case scalingo.DatabaseStatusStopped:
			return retry.NewRetryCancelError(errors.Newf(ctx, "database %s is stopped", appID))
		default:
			return errors.Newf(ctx, "database %s status is %s", appID, db.Database.Status)
		}
	})
	if err != nil {
		return errors.Wrap(ctx, err, "wait for running database")
	}

	return nil
}
