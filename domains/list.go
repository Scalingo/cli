package domains

import (
	"context"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo/v7"
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
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	domains, err := c.DomainsList(ctx, app)
	if err != nil {
		return errgo.Mask(err)
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"Domain", "TLS/SSL", "TLS Subject", "Let's Encrypt Certificate"})
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
			// So if a customer certification is used, it is still mentionned we have it
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
		t.Append(row)
	}
	t.Render()

	if hasCanonical {
		fmt.Println("  (*) canonical domain")
	}
	return nil
}
