package region_migrations

import (
	"fmt"

	"gopkg.in/AlecAivazis/survey.v1"

	scalingo "github.com/Scalingo/go-scalingo/v4"
)

func ConfirmPrepare(migration scalingo.RegionMigration) bool {
	fmt.Printf("Note: Prepare will create an empty canvas. This won't modify your production application\n\n")
	fmt.Println("The following operations will be achieved:")
	fmt.Println(" - Mark your app as migrating (preventing you from modifying it)")
	fmt.Printf(" - Create the new app named %s in the region '%s'\n", migration.DstAppName, migration.Destination)
	fmt.Println(" - Import the last deployment")
	fmt.Println(" - Import environment variables")
	fmt.Println(" - Import SCM configuration")
	fmt.Println(" - Import collaborators")
	fmt.Println(" - Import domains and TLS certificates")
	fmt.Println(" - Import application settings")
	fmt.Println(" - Import application formation")
	fmt.Println(" - Import notifiers")
	fmt.Println(" - Import autoscalers")
	fmt.Println(" - Import alerts")

	return askContinue("Continue?")
}

func ConfirmData(migration scalingo.RegionMigration) bool {
	fmt.Println("The following operations will be achieved:")
	fmt.Println(" - Stop your old app")
	fmt.Printf(" - Create addons on the '%s' region\n", migration.Destination)
	fmt.Println(" - Import addons data")

	return askContinue("Continue?")
}

func ConfirmFinalize(migration scalingo.RegionMigration) bool {
	fmt.Println("The following operations will be achieved:")
	if migration.Status != scalingo.RegionMigrationStatusDataMigrated {
		fmt.Println(" - Stop your old app")
	}
	fmt.Println(" - Start the new app")
	fmt.Printf(" - Redirect the traffic coming to '%s' on the old region to '%s' on '%s'\n", migration.SrcAppName, migration.DstAppName, migration.Destination)

	return askContinue("Continue?")
}

func askContinue(message string) bool {
	result := false
	prompt := &survey.Confirm{
		Message: "Continue?",
	}
	survey.AskOne(prompt, &result, nil)
	return result
}

func ConfirmStep(migration scalingo.RegionMigration, step scalingo.RegionMigrationStep) bool {
	switch step {
	case scalingo.RegionMigrationStepPrepare:
		return ConfirmPrepare(migration)
	case scalingo.RegionMigrationStepData:
		return ConfirmData(migration)
	case scalingo.RegionMigrationStepFinalize:
		return ConfirmFinalize(migration)
	}
	return true
}
