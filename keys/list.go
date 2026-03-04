package keys

import (
	"context"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/go-utils/errors/v3"

	"github.com/Scalingo/cli/config"
)

func List(ctx context.Context) error {
	c, err := config.ScalingoAuthClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}
	keys, err := c.KeysList(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to list SSH keys")
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.Header([]string{"Name", "Content"})

	for _, k := range keys {
		t.Append([]string{k.Name, k.Content[0:20] + "..." + k.Content[len(k.Content)-30:]})
	}

	t.Render()
	return nil
}
