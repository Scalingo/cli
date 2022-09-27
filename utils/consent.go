package utils

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
)

// CheckForConsent will check if an operator does have consent before executing a command. If they doesn't or if the CLI can't determine if the user has consent, it will ask for the operator confirmation
// All display takes place on Stderr to minimize the chance that it will collide in situation where stdout is piped to another process (typically `scalingo logs | grep SOMETHING`)
func CheckForConsent(ctx context.Context, appName string) {
	currentUser, err := config.C.CurrentUser()
	if err != nil {
		return
	}

	// If the user is not admin, exit immediatly, this will make this function
	// almost a NOOP for non operators.
	if !currentUser.Flags["admin"] {
		return
	}

	// From this point out, if we encounter an error, we try to safely recover by manually asking the operator to override.
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		askForConsent(false)
		return
	}

	// Check if the operator is a collaborator on the targetted app
	apps, err := c.AppsList(ctx)
	if err != nil {
		askForConsent(false)
		return
	}

	for _, app := range apps {
		if app.Name == appName {
			// The operator is a collaborator, no consent needed
			return
		}
	}

	// The operator is not a collaborator, checking for consent

	app, err := c.AppsShow(ctx, appName)
	if err != nil {
		askForConsent(false)
		return
	}

	if app.DataAccessConsent == nil {
		// No consent for this app, asking for an override
		askForConsent(true)
		return
	}

	containers := checkAccessContent(app.DataAccessConsent.ContainersUntil)
	databases := checkAccessContent(app.DataAccessConsent.DatabasesUntil)

	if containers && databases {
		// There is a consent for both app and containers, no need for operator validation, it can safely continue
		return
	}

	printAccessContent("Containers Access Content", containers)
	printAccessContent("Database Access Content", databases)

	if !databases && !containers {
		// The operator does not have access to apps nor containers, asking for an override
		askForConsent(true)
		return
	}

	// Here the operator only have access to one of the two consent. We cannot safely assume what the operator will do, asking for manual validation.
	askForConsent(false)
}

func askForConsent(override bool) {
	if override {
		// If the override boolean is set to true, we know that the operator does not have consent for this app, asking for an override
		fmt.Fprint(os.Stderr, io.BoldRed("You do not have consent for this app, Override ? (y/n) "))
	} else {
		// If the override boolean is set to false, we do not know if the operator has consent for this app, asking for confirmation
		fmt.Fprint(os.Stderr, "Do you have consent to access this app? (y/n) ")
	}
	var confirm string
	fmt.Scanln(&confirm)
	if confirm != "y" && confirm != "Y" {
		fmt.Fprintln(os.Stderr, io.BoldRed("No consent given, stopping..."))
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr)
}

func checkAccessContent(t *time.Time) bool {
	value := false
	if t != nil && t.After(time.Now()) {
		value = true
	}

	return value
}

func printAccessContent(message string, value bool) {
	if value {
		fmt.Fprintf(os.Stderr, "%s: %s\n", message, io.Green("true"))
	} else {
		fmt.Fprintf(os.Stderr, "%s: %s\n", message, io.BoldRed("false"))
	}
}
