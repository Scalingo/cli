package privatenetworks

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
)

const containerTypeWeb = "web"

func List(ctx context.Context, app string, format string, page string, perPage string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	domainNames, err := c.PrivateNetworksDomainsList(ctx, app, page, perPage)
	if err != nil {
		return errgo.Notef(err, "fail to list private network domains")
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
				return errgo.Notef(err, "fail to append row to table")
			}
		}

		err := t.Render()
		if err != nil {
			return errgo.Notef(err, "fail to render table")
		}

		fmt.Printf("Page %d/%d\nTotal number of domains %d\n", domainNames.Meta.CurrentPage, domainNames.Meta.TotalPages, domainNames.Meta.TotalCount)
	case "json":
		jsonOutput, err := json.Marshal(domainNames)
		if err != nil {
			return errgo.Notef(err, "fail to generate output JSON from domain names")
		}
		fmt.Println(string(jsonOutput))
	}

	return nil
}
