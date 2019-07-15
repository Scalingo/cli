package region_migrations

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/Scalingo/cli/config"
	scalingo "github.com/Scalingo/go-scalingo"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/gosuri/uilive"
	errgo "gopkg.in/errgo.v1"
)

const (
	SUCCESS = "✔"
	ERROR   = "✘"
)

func WatchMigration(client *scalingo.Client, appId, id string) error {
	writer := uilive.New()
	curStep := 0
	migration, err := client.ShowRegionMigration(appId, id)
	if err != nil {
		return errgo.Notef(err, "fail to find migration")
	}
	lock := &sync.Mutex{}

	go func() {
		for {
			lock.Lock()
			migration, err = client.ShowRegionMigration(appId, id)
			lock.Unlock()
			if err != nil ||
				migration.Status == scalingo.RegionMigrationStatusError ||
				migration.Status == scalingo.RegionMigrationStatusPreflightError {
				return
			}
			time.Sleep(1 * time.Second)
		}

	}()

	for {
		lock.Lock()
		writeMigration(writer, migration, curStep)
		if err != nil {
			lock.Unlock()
			return errgo.Notef(err, "fail tor get migration")
		}
		lock.Unlock()
		if migration.Status != scalingo.RegionMigrationStatusScheduled &&
			migration.Status != scalingo.RegionMigrationStatusRunning {
			break
		}
		curStep = (curStep + 1) % 8
		time.Sleep(100 * time.Millisecond)
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

func writeMigration(w *uilive.Writer, migration scalingo.RegionMigration, curStep int) {
	fmt.Fprintf(w, "Migrating app: %s\n", migration.AppName)
	fmt.Fprintf(w.Newline(), "Destination: %s\n", migration.Destination)
	if migration.NewAppID == "" {
		fmt.Fprintf(w.Newline(), "New app id: %s\n", color.BlueString("N/C"))
	} else {
		fmt.Fprintf(w.Newline(), "New app id: %s\n", migration.NewAppID)
	}
	fmt.Fprintf(w.Newline(), "Status: %s\n", formatMigrationStatus(migration.Status))
	if migration.Status == scalingo.RegionMigrationStatusScheduled {
		fmt.Fprintf(w.Newline(), "%s Waiting for the migration to start\n", spinner.CharSets[11][curStep])
	}

	for i, _ := range migration.Steps {
		writeStep(w, migration.Steps[len(migration.Steps)-1-i], curStep)
	}

	w.Flush()
}

func writeStep(w *uilive.Writer, step scalingo.Step, curStep int) {
	result := ""
	switch step.Status {
	case scalingo.StepStatusRunning:
		result = color.BlueString(fmt.Sprintf("%s %s...", spinner.CharSets[11][curStep], step.Name))
	case scalingo.StepStatusDone:
		result = color.GreenString(fmt.Sprintf("%s %s Done!", SUCCESS, step.Name))
	case scalingo.StepStatusError:
		result = color.RedString(fmt.Sprintf("%s %s FAILED!", ERROR, step.Name))
	}
	fmt.Fprintf(w.Newline(), "%s\n", result)
}

func formatMigrationStatus(status scalingo.RegionMigrationStatus) string {
	strStatus := string(status)
	switch status {
	case scalingo.RegionMigrationStatusScheduled:
		return color.BlueString(strStatus)
	case scalingo.RegionMigrationStatusRunning:
		return color.YellowString(strStatus)
	case scalingo.RegionMigrationStatusDone:
		return color.GreenString(strStatus)
	case scalingo.RegionMigrationStatusPreflightError:
		fallthrough
	case scalingo.RegionMigrationStatusError:
		return color.RedString(strStatus)
	}

	return color.BlueString(strStatus)
}
