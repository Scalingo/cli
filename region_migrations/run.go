package region_migrations

import (
	"context"
	"fmt"

	errgo "gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	scalingo "github.com/Scalingo/go-scalingo/v6"
)

func Create(ctx context.Context, app string, destination string, dstAppName string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get scalingo client")
	}

	migration, err := c.CreateRegionMigration(ctx, app, scalingo.RegionMigrationParams{
		Destination: destination,
		DstAppName:  dstAppName,
	})
	if err != nil {
		return errgo.Notef(err, "fail to create migration")
	}

	err = c.RunRegionMigrationStep(ctx, app, migration.ID, scalingo.RegionMigrationStepPreflight)
	if err != nil {
		return errgo.Notef(err, "fail to run preflight step")
	}

	err = WatchMigration(ctx, c, app, migration.ID, RefreshOpts{
		ExpectedStatuses: []scalingo.RegionMigrationStatus{
			scalingo.RegionMigrationStatusPreflightError,
			scalingo.RegionMigrationStatusPreflightSuccess,
		},
	})
	if err != nil {
		return errgo.Notef(err, "fail to watch migration")
	}

	return nil
}

func Run(ctx context.Context, app, migrationID string, step scalingo.RegionMigrationStep) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get scalingo client")
	}

	migration, err := c.ShowRegionMigration(ctx, app, migrationID)
	if err != nil {
		return errgo.Notef(err, "fail to show region migration")
	}

	shouldContinue := ConfirmStep(migration, step)
	if !shouldContinue {
		fmt.Println("The current step has been canceled. You can restart it later.")
		fmt.Println("If you want to abort the migration, run:")
		fmt.Printf("scalingo --region %s --app %s migration-abort %s\n", migration.Source, app, migrationID)
		return nil
	}
	expectedStatuses := []scalingo.RegionMigrationStatus{}

	switch step {
	case scalingo.RegionMigrationStepPrepare:
		expectedStatuses = append(expectedStatuses, scalingo.RegionMigrationStatusPrepared)
	case scalingo.RegionMigrationStepData:
		expectedStatuses = append(expectedStatuses, scalingo.RegionMigrationStatusDataMigrated)
	}

	previousStepIDs := []string{}

	for _, step := range migration.Steps {
		previousStepIDs = append(previousStepIDs, step.ID)
	}

	err = c.RunRegionMigrationStep(ctx, app, migrationID, step)
	if err != nil {
		return errgo.Notef(err, "fail to run %s step", step)
	}

	err = WatchMigration(ctx, c, app, migrationID, RefreshOpts{
		ExpectedStatuses: expectedStatuses,
		HiddenSteps:      previousStepIDs,
		CurrentStep:      step,
	})
	if err != nil {
		return errgo.Notef(err, "fail to watch migration")
	}

	return nil
}

func Abort(ctx context.Context, app, migrationID string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get scalingo client")
	}

	migration, err := c.ShowRegionMigration(ctx, app, migrationID)
	if err != nil {
		return errgo.Notef(err, "fail to show region migration")
	}

	previousStepIDs := []string{}

	for _, step := range migration.Steps {
		previousStepIDs = append(previousStepIDs, step.ID)
	}

	err = c.RunRegionMigrationStep(ctx, app, migrationID, scalingo.RegionMigrationStepAbort)
	if err != nil {
		return errgo.Notef(err, "fail to run abort step")
	}

	err = WatchMigration(ctx, c, app, migrationID, RefreshOpts{
		ExpectedStatuses: []scalingo.RegionMigrationStatus{
			scalingo.RegionMigrationStatusAborted,
		},
		HiddenSteps: previousStepIDs,
		CurrentStep: scalingo.RegionMigrationStepAbort,
	})
	if err != nil {
		return errgo.Notef(err, "fail to watch migration")
	}

	return nil
}
