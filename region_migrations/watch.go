package region_migrations

import (
	"fmt"
	"net"
	"net/url"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/utils"
	scalingo "github.com/Scalingo/go-scalingo"
	"github.com/fatih/color"
	errgo "gopkg.in/errgo.v1"
)

const (
	SUCCESS = "✔"
	ERROR   = "✘"
)

func WatchMigration(client *scalingo.Client, appId, id string) error {
	refresher := NewRefresher(client, appId, id)
	migration, err := refresher.Start()
	if err != nil {
		return errgo.Notef(err, "fail to watch migration")
	}

	if migration == nil {
		return nil
	}

	if migration.Status == scalingo.RegionMigrationStatusDone {
		newRegionClient, err := config.ScalingoClientForRegion(migration.Destination)
		if err != nil {
			showGenericMigrationSuccessMessage()
			return nil
		}

		app, err := newRegionClient.AppsShow(appId)
		if err != nil {
			showGenericMigrationSuccessMessage()
			return nil
		}

		domains, err := newRegionClient.DomainsList(appId)
		if err != nil {
			showGenericMigrationSuccessMessage()
			return nil
		}

		color.Green("\n\n\nYour application is now available at: %s\n\n", app.BaseURL)

		if len(domains) == 0 {
			return nil
		}

		dnsARecord := ""
		dnsCNAMERecord := ""
		parsed, err := url.Parse(app.BaseURL)
		if err != nil {
			io.Warning("You need to change the DNS record of your domains to point to the new region.")
			io.Status("See: https://doc.scalingo.com/platform/app/domain#configure-your-domain-name for more informations")
			return nil
		}
		dnsCNAMERecord = parsed.Host
		addresses, err := net.LookupHost(parsed.Host)
		if err == nil && len(addresses) > 0 {
			dnsARecord = addresses[0]
		}

		showDoc := false
		io.Status("You need to make the following changes to your DNS records:\n")
		for _, domain := range domains {
			isCNAME, err := utils.IsCNAME(domain.Name)
			if err != nil || (!isCNAME && dnsARecord == "") {
				io.Infof("- %s record should be changed\n", domain.Name)
				showDoc = true
			}
			if isCNAME {
				io.Infof("- CNAME record of %s should be changed to %s\n", domain.Name, dnsCNAMERecord)
			} else {
				io.Infof("- A record of %s should be changed to %s\n", domain.Name, dnsARecord)
			}
		}

		if showDoc {
			io.Info("To configure your DNS record, check out our documentation:")
			io.Info("- https://doc.scalingo.com/platform/app/domain#configure-your-domain-name")
		}
		return nil
	}

	color.Red("The migration failed because of the following errors:\n")

	for i, _ := range migration.Steps {
		step := migration.Steps[len(migration.Steps)-1-i]
		if step.Status == scalingo.StepStatusError {
			if step.Logs == "" {
				color.Red("- Step %s failed\n", step.Name)
			} else {
				color.Red("- Step %s failed: %s\n", step.Name, step.Logs)
			}
		}
	}

	fmt.Println("The application has been rolled back to its working state.")
	fmt.Println("Contact support@scalingo.com to troubleshoot this issue.")

	return nil
}

func showGenericMigrationSuccessMessage() {
	color.Green("\n\n\nThe application has been migrated!\n")
	fmt.Println("You need to change the DNS record of your domains to point to the new region.")
	fmt.Println("See: https://doc.scalingo.com/platform/app/domain#configure-your-domain-name for more informations")
}
