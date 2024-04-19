package domains

import (
	"context"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	scalingo "github.com/Scalingo/go-scalingo/v7"
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
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	domains, err := c.DomainsList(ctx, app)
	if err != nil {
		return errgo.Mask(err)
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"Domain", "TLS/SSL"})
	hasCanonical := false

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
			row = append(row, fmt.Sprintf("Let's Encrypt: %s", letsencryptStatus))
		} else {
			row = append(row, fmt.Sprintf("Valid until %v", domain.Validity))
		}
		t.Append(row)
	}
	t.Render()

	if hasCanonical {
		fmt.Println("  (*) canonical domain")
	}
	return nil
}
