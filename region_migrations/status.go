package region_migrations

import (
	"fmt"
	"net"
	"net/url"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
	scalingo "github.com/Scalingo/go-scalingo"
	"github.com/fatih/color"
)

func showMigrationStatusSuccess(appId string, migration scalingo.RegionMigration) {
	newRegionClient, err := config.ScalingoClientForRegion(migration.Destination)
	if err != nil {
		showGenericMigrationSuccessMessage()
		return
	}

	app, err := newRegionClient.AppsShow(appId)
	if err != nil {
		showGenericMigrationSuccessMessage()
		return
	}

	domains, err := newRegionClient.DomainsList(appId)
	if err != nil {
		showGenericMigrationSuccessMessage()
		return
	}

	color.Green("Your application is now available at: %s\n\n", app.BaseURL)

	if len(domains) == 0 {
		return
	}

	dnsARecord := ""
	dnsCNAMERecord := ""
	parsed, err := url.Parse(app.BaseURL)
	if err != nil {
		io.Warning("You need to change the DNS record of your domains to point to the new region.")
		io.Status("See: https://doc.scalingo.com/platform/app/domain#configure-your-domain-name for more informations")
		return
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
	return
}
func showGenericMigrationSuccessMessage() {
	color.Green("The application has been migrated!\n")
	fmt.Println("You need to change the DNS record of your domains to point to the new region.")
	fmt.Println("See: https://doc.scalingo.com/platform/app/domain#configure-your-domain-name for more information")
}

func showMigrationStatusFailed(appId string, migration scalingo.RegionMigration, opts RefreshOpts) {
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

	if opts.CurrentStep != scalingo.RegionMigrationStepAbort {
		fmt.Println("To rollback your application to a working state, run:")
		fmt.Printf("scalingo --app %s migration-abort %s\n\n", appId, migration.ID)
	}

	fmt.Println("You can contact support@scalingo.com to troubleshoot this issue.")
}

func showMigrationStatusPreflightSuccess(appId string, migration scalingo.RegionMigration) {
	fmt.Printf("Your app can be migrated to the %s zone.\n", migration.Destination)
	fmt.Printf("To start the migration launch:\n\n")
	fmt.Printf("scalingo --app %s migration-run --prepare %s\n", appId, migration.ID)
}

func showMigrationStatusPrepared(appId string, migration scalingo.RegionMigration) {
	fmt.Printf("Application on region '%s' has been prepared, you can now:\n", migration.Destination)
	fmt.Printf("- Let us migrate your data to '%s' newly created databases with:\n", migration.Destination)
	fmt.Printf("scalingo --app %s migration-run --data %s\n", appId, migration.ID)
	fmt.Printf("- Handle data migration manually, then finalizing the migration with:\n")
	fmt.Printf("scalingo --app %s migration-run --finalize %s\n", appId, migration.ID)
}

func showMigrationStatusDataMigrated(appId string, migration scalingo.RegionMigration) {
	fmt.Printf("Data has been migrated to the '%s' region\n", migration.Destination)
	fmt.Printf("You can finalize the migration with:\n")
	fmt.Printf("scalingo --app %s migration-run --finalize %s\n", appId, migration.ID)
}

func showMigrationStatusAborted(appId string, migration scalingo.RegionMigration) {
	fmt.Printf("The migration '%s' has been aborted\n", migration.ID)
	fmt.Printf("You can retry it with:\n")
	fmt.Printf("scalingo --app %s migration-create --to %s", appId, migration.Destination)
	if migration.DstAppName != migration.SrcAppName {
		fmt.Printf(" --new-name %s \n", migration.DstAppName)
	} else {
		fmt.Printf("\n")
	}
}
