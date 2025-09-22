package privatenetworks

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-utils/errors/v2"
)

const containerTypeWeb = "web"

func List(ctx context.Context, app string, format string, page uint, perPage uint) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "get Scalingo client")
	}

	domainNames, err := c.PrivateNetworksDomainsList(ctx, app, page, perPage)
	if err != nil {
		return errors.Wrapf(ctx, err, "list private network domains")
	}

	switch format {
	case "table":
		t := tablewriter.NewWriter(os.Stdout)
		t.Header([]string{"Container Type", "Domain Name"})

		for _, domain := range domainNames.Data {
			containerType := containerTypeWeb
			parts := strings.Split(domain, ".")
			slices.Reverse(parts)
			if len(parts) == 5 {
				containerType = containerTypeWeb
			} else if len(parts) > 5 {
				containerType = parts[5]
			}
			err := t.Append([]string{containerType, domain})
			if err != nil {
				return errors.Wrapf(ctx, err, "append row to table")
			}
		}

		err := t.Render()
		if err != nil {
			return errors.Wrapf(ctx, err, "render table")
		}

		fmt.Printf("Page %d/%d\nTotal number of domains %d\n", domainNames.Meta.CurrentPage, domainNames.Meta.TotalPages, domainNames.Meta.TotalCount)
	case "json":
		jsonOutput, err := json.Marshal(domainNames)
		if err != nil {
			return errors.Wrapf(ctx, err, "generate output JSON from domain names")
		}
		fmt.Println(string(jsonOutput))
	}

	return nil
}
