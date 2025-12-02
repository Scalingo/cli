package projects

import (
	"context"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/go-utils/logger"
)

func Get(ctx context.Context, projectID string) error {
	log := logger.Get(ctx)
	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	project, err := client.ProjectGet(ctx, projectID)
	if err != nil {
		return errors.Wrap(ctx, err, "get project")
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.Header([]string{"Project Field", "Value"})

	_ = t.Append([]string{"Name", project.Name})
	_ = t.Append([]string{"ID", project.ID})
	_ = t.Append([]string{"Default", strconv.FormatBool(project.Default)})
	_ = t.Append([]string{"Owner", project.Owner.Username})

	if project.Flags["private-network"] {
		_ = t.Append([]string{"", ""})
		_ = t.Append([]string{"Private Network", "true"})

		privateNetworkInfo, err := client.ProjectPrivateNetworkGet(ctx, projectID)
		if err != nil {
			log.WithError(err).Error("Failed to fetch private network info")
			_ = t.Append([]string{"", "Failed to fetch private network info"})
		} else {
			_ = t.Append([]string{" - ID", privateNetworkInfo.ID})
			_ = t.Append([]string{" - Subnet", privateNetworkInfo.Subnet})
			_ = t.Append([]string{" - Gateway IP", privateNetworkInfo.Gateway})
			_ = t.Append([]string{" - Total number of assignable IPs", strconv.Itoa(privateNetworkInfo.MaxIPsCount)})
			_ = t.Append([]string{" - Used IPs count", strconv.Itoa(privateNetworkInfo.UsedIPsCount)})

			if len(privateNetworkInfo.UsedIPs) == 0 {
				_ = t.Append([]string{" - Used IPs", "None"})
			} else {
				_ = t.Append([]string{" - Used IPs", privateNetworkInfo.UsedIPs[0]})
			}
			for _, usedIP := range privateNetworkInfo.UsedIPs[1:] {
				_ = t.Append([]string{"", usedIP})
			}
		}
	}

	_ = t.Render()

	return nil
}
