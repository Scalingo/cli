package privatenetworks

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"sort"
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

	if len(domainNames) == 0 {
		fmt.Println("No private network domain found")
		return nil
	}

	switch format {
	case "table":
		sort.Slice(domainNames, func(i, j int) bool {
			iParts := strings.Split(string(domainNames[i]), ".")
			jParts := strings.Split(string(domainNames[j]), ".")

			// compare from the end of the domain name
			for i := 2; i <= len(iParts) && i <= len(jParts); i++ {
				if iParts[len(iParts)-i] == jParts[len(jParts)-i] {
					continue
				}

				// place "web" container types at the top of the list
				if i == 5 {
					if iParts[len(iParts)-i] == containerTypeWeb && jParts[len(jParts)-i] != containerTypeWeb {
						return true
					}
					if iParts[len(iParts)-i] != containerTypeWeb && jParts[len(jParts)-i] == containerTypeWeb {
						return false
					}
				}
				return iParts[len(iParts)-i] < jParts[len(jParts)-i]
			}

			// shorter domain names should be at the top of the list
			return len(iParts) < len(jParts)
		})

		t := tablewriter.NewWriter(os.Stdout)
		t.Header([]string{"Container Type", "Domain Name"})

		for _, domain := range domainNames {
			containerType := containerTypeWeb
			parts := strings.Split(string(domain), ".")
			slices.Reverse(parts)
			if len(parts) == 4 {
				containerType = containerTypeWeb
			} else if len(parts) > 4 {
				containerType = parts[4]
			}
			err := t.Append([]string{containerType, string(domain)})
			if err != nil {
				return errgo.Notef(err, "fail to append row to table")
			}
		}
		err := t.Render()
		if err != nil {
			return errgo.Notef(err, "fail to render table")
		}
	case "json":
		jsonOutput, err := json.Marshal(domainNames)
		if err != nil {
			return errgo.Notef(err, "fail to generate output JSON from domain names")
		}
		fmt.Println(string(jsonOutput))
	}

	return nil
}
