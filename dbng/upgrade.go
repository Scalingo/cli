package dbng

import (
	"context"
	"time"

	"github.com/briandowns/spinner"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v9"
	"github.com/Scalingo/go-utils/errors/v2"
)

func Upgrade(ctx context.Context, databaseID, plan string, wait bool) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	db, err := c.Preview().DatabaseShow(ctx, databaseID)
	if err != nil {
		return errors.Wrap(ctx, err, "show database")
	}

	planID, err := utils.FindPlan(ctx, c, db.Technology, plan)
	if err != nil {
		return errors.Wrap(ctx, err, "invalid database plan")
	}

	addons, err := c.AddonsList(ctx, db.App.ID)
	if err != nil {
		return errors.Wrap(ctx, err, "list database addons")
	}
	if len(addons) == 0 {
		return errors.Newf(ctx, "no addons found for database application %s", db.App.ID)
	}

	var dbAddon *scalingo.Addon
	if len(addons) == 1 {
		dbAddon = addons[0]
	} else if db.Database.ResourceID != "" {
		for _, addon := range addons {
			if addon.ResourceID == db.Database.ResourceID {
				dbAddon = addon
				break
			}
		}
	}
	if dbAddon == nil {
		return errors.Newf(ctx, "no addon matching database %s (%s) found for application %s", db.Name, db.ID, db.App.ID)
	}

	_, err = c.AddonUpgrade(ctx, db.App.ID, dbAddon.ID, scalingo.AddonUpgradeParams{
		PlanID: planID,
	})
	if err != nil {
		return errors.Wrap(ctx, err, "upgrade database")
	}

	io.Statusf("Your %s database %s ('%s') is being upgraded…\n\n",
		db.Technology,
		db.ID,
		db.Name,
	)

	if wait {
		spin := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		spin.Suffix = " Waiting for the plan change to complete\n"
		spin.Start()
		defer spin.Stop()

		err = waitForDatabasePlanChange(ctx, c, db.ID)
		if err != nil {
			return errors.Wrap(ctx, err, "wait for database plan change")
		}

		io.Status("Your database plan change is complete.")
	}

	return nil
}

func waitForDatabasePlanChange(ctx context.Context, client *scalingo.Client, databaseID string) error {
	db, err := client.Preview().DatabaseShow(ctx, databaseID)
	if err != nil {
		return err
	}

	if db.Database.Status == scalingo.DatabaseStatusRunning {
		err = waitForDatabaseStatus(ctx, client, databaseID, func(status scalingo.DatabaseStatus) bool {
			return status != scalingo.DatabaseStatusRunning
		})
		if err != nil {
			return err
		}
	}

	return waitForRunningDatabase(ctx, client, databaseID)
}

func waitForDatabaseStatus(ctx context.Context, client *scalingo.Client, databaseID string, check func(scalingo.DatabaseStatus) bool) error {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		db, err := client.Preview().DatabaseShow(ctx, databaseID)
		if errors.Is(err, scalingo.ErrDatabaseNotFound) {
			continue
		} else if err != nil {
			return err
		}
		if check(db.Database.Status) {
			break
		}
	}
	return nil
}
