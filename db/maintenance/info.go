package maintenance

import (
	"context"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-utils/errors/v2"
)

func Info(ctx context.Context, app, addonUUID, maintenanceID string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Notef(ctx, err, "get Scalingo client")
	}

	maintenanceInfo, err := c.DatabaseShowMaintenance(ctx, app, addonUUID, maintenanceID)
	if err != nil {
		return errors.Notef(ctx, err, "get the maintenance")
	}

	var startedAt string
	var endedAt string

	if maintenanceInfo.StartedAt != nil {
		startedAt = maintenanceInfo.StartedAt.Local().Format(utils.TimeFormat)
	}
	if maintenanceInfo.EndedAt != nil {
		endedAt = maintenanceInfo.EndedAt.Local().Format(utils.TimeFormat)
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.Append([]string{"ID", fmt.Sprintf("%v", maintenanceInfo.ID)})
	t.Append([]string{"Type", fmt.Sprintf("%v", maintenanceInfo.Type)})
	t.Append([]string{"Started At", fmt.Sprintf("%v", startedAt)})
	t.Append([]string{"Ended At", fmt.Sprintf("%v", endedAt)})
	t.Append([]string{"Status", fmt.Sprintf("%v", maintenanceInfo.Status)})

	t.Render()

	return nil
}
