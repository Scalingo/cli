package regions

import (
	"context"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-utils/errors/v3"
)

func List(ctx context.Context) error {
	c, err := config.ScalingoAuthClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get scalingo API client")
	}
	regions, err := c.RegionsList(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to list available regions")
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.Header([]string{"Name", "Display", "API Endpoint"})

	for _, r := range regions {
		t.Append([]string{r.Name, r.DisplayName, r.API})
	}

	t.Render()
	return nil
}
