package region_migrations

import (
	"fmt"
	"net/url"

	"github.com/Scalingo/cli/config"
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
			color.Green("The application has been migrated!\n")
			return nil
		}

		app, err := newRegionClient.AppsShow(appId)
		if err != nil {
			color.Green("The application has been migrated!\n")
			return nil
		}

		domains, err := newRegionClient.DomainsList(appId)
		if err != nil {
			color.Green("The application has been migrated!\n")
			return nil
		}

		color.Green("Your application is now available at: %s\n\n", app.BaseURL)
		if len(domains) > 0 {
			parsed, err := url.Parse(app.BaseURL)
			if err != nil {
				fmt.Printf("You need to change the CNAME record of the yours domains to point to the new region.")
				return nil
			}
			fmt.Printf("You need to change the CNAME record of the following domains to point to '%s'\n", parsed.Host)
			for _, domain := range domains {
				fmt.Printf(" - %s\n", domain.Name)
			}
		}
		return nil
	}

	color.Red("The migration failed because of the following error:\n")

	for i, _ := range migration.Steps {
		step := migration.Steps[len(migration.Steps)-1-i]
		if step.Status == scalingo.StepStatusError {
			if step.Logs == "" {
				color.Red("- The step: %s failed.\n", step.Name)
			} else {
				color.Red("- Step %s failed: %s\n", step.Name, step.Logs)
			}
		}
	}

	fmt.Println("The application has been rolled back to it's working state.")
	fmt.Println("Contact support@scalingo.com to troubleshoot this issue.")

	return nil
}
