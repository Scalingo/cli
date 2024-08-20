package regionmigrations

import (
	"github.com/fatih/color"

	"github.com/Scalingo/go-scalingo/v7"
)

func formatMigrationStatus(status scalingo.RegionMigrationStatus) string {
	strStatus := string(status)
	switch status {
	case scalingo.RegionMigrationStatusPrepared:
		fallthrough
	case scalingo.RegionMigrationStatusDataMigrated:
		fallthrough
	case scalingo.RegionMigrationStatusAborting:
		fallthrough
	case scalingo.RegionMigrationStatusCreated:
		return color.BlueString(strStatus)
	case scalingo.RegionMigrationStatusRunning:
		return color.YellowString(strStatus)
	case scalingo.RegionMigrationStatusPreflightSuccess:
		fallthrough
	case scalingo.RegionMigrationStatusDone:
		return color.GreenString(strStatus)
	case scalingo.RegionMigrationStatusPreflightError:
		fallthrough
	case scalingo.RegionMigrationStatusAborted:
		fallthrough
	case scalingo.RegionMigrationStatusError:
		return color.RedString(strStatus)
	}

	return color.BlueString(strStatus)
}
