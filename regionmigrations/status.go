package regionmigrations

import (
	"context"
	"fmt"
	"net"
	"net/url"

	"github.com/fatih/color"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v7"
)

func showMigrationStatusSuccess(ctx context.Context, appID string, migration scalingo.RegionMigration) {
	newRegionClient, err := config.ScalingoClientForRegion(ctx, migration.Destination)
	if err != nil {
		showGenericMigrationSuccessMessage()
		return
	}

	app, err := newRegionClient.AppsShow(ctx, appID)
	if err != nil {
		showGenericMigrationSuccessMessage()
		return
	}

	domains, err := newRegionClient.DomainsList(ctx, appID)
	if err != nil {
		showGenericMigrationSuccessMessage()
		return
	}

	color.Green("Your application is now available at: %s\n\n", app.BaseURL)

	io.Status("You will also need to change the Git URL of your repository.")
	io.Info("To change the Git remote URL use:")
	io.Infof("git remote set-url scalingo %s \n\n", app.GitURL)

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
}
func showGenericMigrationSuccessMessage() {
	color.Green("The application has been migrated!\n")
	fmt.Println("You need to change the DNS record of your domains to point to the new region.")
	fmt.Println("See: https://doc.scalingo.com/platform/app/domain#configure-your-domain-name for more information")
}

func showMigrationStatusFailed(appID string, migration scalingo.RegionMigration, opts RefreshOpts) {
	color.Red("The migration failed because of the following errors:\n")

	for i := range migration.Steps {
		step := migration.Steps[len(migration.Steps)-1-i]
		if step.Status == scalingo.StepStatusError {
			if step.Logs == "" {
				color.Red("- Step %s failed\n", step.Name)
			} else {
				color.Red("- Step %s failed: %s\n", step.Name, step.Logs)
			}
		}
	}

	if opts.CurrentStep != scalingo.RegionMigrationStepAbort && migration.Status != scalingo.RegionMigrationStatusPreflightError {
		fmt.Println("To rollback your application to a working state, run:")
		fmt.Printf("scalingo --region %s --app %s migration-abort %s\n\n", migration.Source, appID, migration.ID)
	}

	fmt.Println("You can contact support@scalingo.com to troubleshoot this issue.")
}

func showMigrationStatusPreflightSuccess(appID string, migration scalingo.RegionMigration) {
	fmt.Printf("Your app can be migrated to the %s zone.\n\n", migration.Destination)
	fmt.Printf("- Start the migration with:\n")
	fmt.Printf("scalingo --region %s --app %s migration-run --prepare %s\n", migration.Source, appID, migration.ID)
	fmt.Printf("- Abort the migration with:\n")
	fmt.Printf("scalingo --region %s --app %s migration-abort %s\n", migration.Source, appID, migration.ID)
}

func showMigrationStatusPrepared(appID string, migration scalingo.RegionMigration) {
	fmt.Printf("Application on region '%s' has been prepared, you can now:\n", migration.Destination)
	fmt.Printf("- Let us migrate your data to '%s' newly created databases with:\n", migration.Destination)
	fmt.Printf("scalingo --region %s --app %s migration-run --data %s\n", migration.Source, appID, migration.ID)
	fmt.Printf("- Handle data migration manually, then finalizing the migration with:\n")
	fmt.Printf("scalingo --region %s --app %s migration-run --finalize %s\n", migration.Source, appID, migration.ID)
}

func showMigrationStatusDataMigrated(appID string, migration scalingo.RegionMigration) {
	fmt.Printf("Data has been migrated to the '%s' region\n", migration.Destination)
	fmt.Printf("You can finalize the migration with:\n")
	fmt.Printf("scalingo --region %s --app %s migration-run --finalize %s\n", migration.Source, appID, migration.ID)
}

func showMigrationStatusAborted(appID string, migration scalingo.RegionMigration) {
	fmt.Printf("The migration '%s' has been aborted\n", migration.ID)
	fmt.Printf("You can retry it with:\n")
	fmt.Printf("scalingo --region %s --app %s migration-create --to %s", migration.Source, appID, migration.Destination)
	if migration.DstAppName != migration.SrcAppName {
		fmt.Printf(" --new-name %s \n", migration.DstAppName)
	} else {
		fmt.Printf("\n")
	}
}
