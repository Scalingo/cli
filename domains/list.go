package domains

import (
	"context"
	"fmt"
	"os"
	"slices"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo/v9"
	"github.com/Scalingo/go-utils/errors/v2"
)

var letsencryptStatusString = map[string]string{
	string(scalingo.LetsEncryptStatusPendingDNS):  "Pending DNS",
	string(scalingo.LetsEncryptStatusNew):         "Creating",
	string(scalingo.LetsEncryptStatusCreated):     "Created",
	string(scalingo.LetsEncryptStatusDNSRequired): "DNS required",
	string(scalingo.LetsEncryptStatusError):       "Error",
}

func List(ctx context.Context, app string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client to list domains")
	}
	domains, err := c.DomainsList(ctx, app)
	if err != nil {
		return errors.Wrap(ctx, err, "list domains")
	}

	t := tablewriter.NewWriter(os.Stdout)
	headers := []string{"Domain", "TLS/SSL", "TLS Subject", "Let's Encrypt Certificate"}
	hasCanonical := false

	for _, domain := range domains {
		domainName := domain.Name
		if domain.Canonical {
			hasCanonical = true
			domainName += " (*)"
		}

		tls := "-"
		letsEncrypt := "Disabled"
		if domain.LetsEncryptEnabled {
			// If the domain is using Let's Encrypt (and not a custom cert), we mention it
			if domain.LetsEncrypt {
				tls = "Let's Encrypt"
			}
			// In any case we display the state of creation of the Let's Encrypt certificate
			// So if a customer certification is used, it is still mentioned we have it
			var ok bool
			letsEncrypt, ok = letsencryptStatusString[string(domain.LetsEncryptStatus)]
			if !ok {
				letsEncrypt = string(domain.LetsEncryptStatus)
			}

			if !domain.LetsEncrypt && domain.LetsEncryptStatus == scalingo.LetsEncryptStatusCreated {
				letsEncrypt = "Created, Not in use"
			}
		}

		if domain.SSL && !domain.LetsEncrypt {
			tls = fmt.Sprintf("Valid until %v", domain.Validity)
		}

		row := []string{domainName, tls, domain.TLSCert, letsEncrypt}

		if domain.LetsEncryptStatus == scalingo.LetsEncryptStatusDNSRequired {
			if !slices.Contains(headers, "Manual Action") {
				headers = append(headers, "Manual Action")
			}
			row = append(row, fmt.Sprintf("%v \n%v", domain.AcmeDNSFqdn, domain.AcmeDNSValue))
		}

		t.Append(row)
	}

	t.Header(headers)
	t.Render()

	if hasCanonical {
		fmt.Println("  (*) canonical domain")
	}
	return nil
}
