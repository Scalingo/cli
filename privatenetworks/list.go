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
	scalingoClient, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "get Scalingo client")
	}

	domainNames, err := scalingoClient.PrivateNetworksDomainsList(ctx, app, page, perPage)
	if err != nil {
		return errors.Wrapf(ctx, err, "list private network domains")
	}

	switch format {
	case "table":
		t := tablewriter.NewWriter(os.Stdout)
		t.Header([]string{"Container Type", "Domain Name"})

		// The container type is the 6th part of the domain name when split by "." and reversed
		// e.g. for "1.web.ap-9142978a-3f9d-48d3-8caa-b042d401ac30.pn-ad0fd6a1-d05e-40ea-bf63-c4f8a75a9d8c.private-network.internal.",
		// the container type is "web" when counting from the right (0-based index)
		const containerTypeIndex = 5
		for _, domain := range domainNames.Data {
			containerType := containerTypeWeb
			parts := strings.Split(domain, ".")
			slices.Reverse(parts)
			if len(parts) == containerTypeIndex {
				containerType = containerTypeWeb
			} else if len(parts) > containerTypeIndex {
				containerType = parts[containerTypeIndex]
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
		jsonOutput, _ := json.Marshal(domainNames)
		fmt.Println(string(jsonOutput))
	}

	return nil
}
