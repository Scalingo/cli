package domains

import (
	"fmt"
	"os"

	"github.com/Scalingo/cli/config"
	scalingo "github.com/Scalingo/go-scalingo"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"
)

var letsencryptStatusString = map[string]string{
	string(scalingo.LetsEncryptStatusPendingDNS):  "Pending DNS",
	string(scalingo.LetsEncryptStatusNew):         "Creating",
	string(scalingo.LetsEncryptStatusCreated):     "In use",
	string(scalingo.LetsEncryptStatusDNSRequired): "DNS required",
	string(scalingo.LetsEncryptStatusError):       "Error",
}

func List(app string) error {
	c := config.ScalingoClient()
	domains, err := c.DomainsList(app)
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
		} else {
			var tls string
			if domain.LetsEncrypt {
				letsencryptStatus, ok := letsencryptStatusString[string(domain.LetsEncryptStatus)]
				if !ok {
					letsencryptStatus = string(domain.LetsEncryptStatus)
				}
				tls = fmt.Sprintf("(LE %s)", letsencryptStatus)
			}
			row = append(row, fmt.Sprintf("%s Valid until %v", tls, domain.Validity))
		}
		t.Append(row)
	}
	t.Render()

	if hasCanonical {
		fmt.Println("  (*) canonical domain")
	}
	return nil
}
