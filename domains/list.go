package domains

import (
	"context"
	"fmt"
	"os"
	"slices"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/cli/config"
	scalingo "github.com/Scalingo/go-scalingo/v6"
	"github.com/Scalingo/go-utils/errors/v2"
)

var letsencryptStatusString = map[string]string{
	string(scalingo.LetsEncryptStatusPendingDNS):  "Pending DNS",
	string(scalingo.LetsEncryptStatusNew):         "Creating",
	string(scalingo.LetsEncryptStatusCreated):     "In use",
	string(scalingo.LetsEncryptStatusDNSRequired): "DNS required",
	string(scalingo.LetsEncryptStatusError):       "Error",
}

func List(ctx context.Context, app string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "get Scalingo client to list domains")
	}

	domains, err := c.DomainsList(ctx, app)
	if err != nil {
		return errors.Wrapf(ctx, err, "list domains")
	}

	t := tablewriter.NewWriter(os.Stdout)
	hasCanonical := false
	headers := []string{"Domain", "TLS/SSL"}

	for _, domain := range domains {
		domainName := domain.Name
		if domain.Canonical {
			hasCanonical = true
			domainName += " (*)"
		}
		row := []string{domainName}

		if !domain.SSL {
			row = append(row, "-")
		} else if domain.LetsEncrypt {
			letsencryptStatus, ok := letsencryptStatusString[string(domain.LetsEncryptStatus)]
			if !ok {
				letsencryptStatus = string(domain.LetsEncryptStatus)
			}
			row = append(row, "Let's Encrypt "+letsencryptStatus)
		} else {
			row = append(row, fmt.Sprintf("Valid until %v", domain.Validity))
		}

		if domain.LetsEncryptStatus == scalingo.LetsEncryptStatusDNSRequired {
			if !slices.Contains(headers, "Manual Action") {
				headers = append(headers, "Manual Action")
			}
			row = append(row, fmt.Sprintf("%v \n%v", domain.AcmeDNSFqdn, domain.AcmeDNSValue))
		}

		t.Append(row)
	}

	t.SetHeader(headers)
	t.Render()

	if hasCanonical {
		fmt.Println("  (*) canonical domain")
	}
	return nil
}
