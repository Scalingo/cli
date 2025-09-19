package detect

import (
	"context"
	"fmt"
	"os"
	"strings"

	stderrors "github.com/pkg/errors"
	"github.com/urfave/cli/v3"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v8/debug"
	"github.com/Scalingo/go-utils/errors/v2"
)

var errDatabaseNotFound = stderrors.New("database name not found")

// GetCurrentResource is the new helper to get the current resource (app or database).
func GetCurrentResource(ctx context.Context, c *cli.Command) string {
	resource, _ := GetCurrentResourceAndDatabase(ctx, c)
	return resource
}

// GetCurrentResourceAndDatabase returns the current resource (app or database)
// and the current database UUID if any.
// It exits CLI in case of error.
func GetCurrentResourceAndDatabase(ctx context.Context, c *cli.Command) (string, string) {
	if os.Getenv("SCALINGO_PREVIEW_FEATURES") != "true" {
		return CurrentApp(c), ""
	}

	currentApp := extractAppName(c)
	currentDatabase, databaseUUID, err := currentDatabaseNameAndUUID(ctx, c)
	if err != nil && !errors.Is(err, errDatabaseNotFound) {
		fmt.Println(err)
		os.Exit(1)
	}

	// Check if --app flag is set explicitly
	var appFlagSet bool
	for _, cliContext := range c.Lineage() {
		if cliContext.IsSet("app") {
			appFlagSet = true
			break
		}
	}

	if appFlagSet && currentDatabase != "" && databaseUUID != "" {
		fmt.Println("You can't use --app and --database flags together.")
		os.Exit(1)
	}

	if currentApp == "" && currentDatabase == "" {
		fmt.Println("No application or database found. Please use --app or --database flag.")
		os.Exit(1)
	}

	currentResource := currentApp
	if currentDatabase != "" {
		currentResource = currentDatabase
	}

	debug.Println("[detect] Current resource is", currentResource)
	if databaseUUID != "" {
		debug.Println("[detect] Current database UUID is", databaseUUID)
	}

	return currentResource, databaseUUID
}

// CurrentApp returns the app name if it has been found in one of the following:
// app flag, environment variable "SCALINGO_APP", current Git remote.
// It returns an empty string if not found.
func CurrentApp(c *cli.Command) string {
	appName := extractAppName(c)

	if appName == "" {
		fmt.Println("Unable to find the application name, please use --app flag.")
		os.Exit(1)
	}
	debug.Println("[detect] App name is", appName)

	return appName
}

// currentDatabaseNameAndUUID returns the database name and its UUID
// if database name has been found in one of the following:
// database flag, environment variable "SCALINGO_DATABASE".
// It returns an empty string and errDatabaseNotFound error if the database name
// is not provided or resource ID not found.
func currentDatabaseNameAndUUID(ctx context.Context, c *cli.Command) (string, string, error) {
	dbName := extractDatabaseNameFromCommandLineOrEnv(c)
	if dbName == "" {
		return "", "", errDatabaseNotFound
	}
	debug.Println("[detect] Database name is", dbName)

	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return "", "", errors.Wrap(ctx, err, "get Scalingo client")
	}

	// In case of a database, addons list should return only one element.
	addons, err := client.AddonsList(ctx, dbName)
	if err != nil {
		return "", "", errors.Wrap(ctx, err, "list addons")
	}

	if len(addons) == 0 {
		fmt.Println("No database found with the given name.")
		return "", "", errDatabaseNotFound
	}

	databaseUUID := addons[0].ID

	if databaseUUID == "" {
		fmt.Println("Unable to find the database ID with the given name.")
		return "", "", errDatabaseNotFound
	}

	debug.Println("[detect] Database UUID is", databaseUUID)

	return dbName, databaseUUID, nil
}

// GetAppNameFromGitRemote searches into the current directory and its parent for a remote
// named remoteName or scalingo-<remoteName>.
//
// It returns the application name and an error.
func GetAppNameFromGitRemote(directory string, remoteName string) (string, error) {
	remotes, err := utils.ScalingoGitRemotes(directory)
	if err != nil {
		return "", err
	}

	altRemoteName := "scalingo-" + remoteName
	for _, remote := range remotes {
		if remote.Config().Name == remoteName ||
			remote.Config().Name == altRemoteName {
			return extractAppNameFromGitRemote(remote.Config().URLs[0]), nil
		}
	}

	return "", errgo.Newf("[detect] Scalingo Git remote hasn't been found")
}

// RemoteNameFromFlags returns the remote name specified in command flags
func RemoteNameFromFlags(c *cli.Command) string {
	for _, cliContext := range c.Lineage() {
		if cliContext.String("remote") != "" {
			return cliContext.String("remote")
		}
	}
	return ""
}

func extractAppName(c *cli.Command) string {
	for _, cliContext := range c.Lineage() {
		appName := cliContext.String("app")
		if appName != "" && appName != "<name>" {
			return appName
		}
	}

	var err error
	var appName string

	if os.Getenv("SCALINGO_APP") != "" {
		appName = os.Getenv("SCALINGO_APP")
	} else if dir, ok := utils.DetectGit(); ok {
		appName, err = GetAppNameFromGitRemote(dir, RemoteNameFromFlags(c))
		if err != nil {
			debug.Println(err)
		}
	}
	return appName
}

func extractDatabaseNameFromCommandLineOrEnv(c *cli.Command) string {
	if os.Getenv("SCALINGO_PREVIEW_FEATURES") != "true" {
		return ""
	}

	for _, cliContext := range c.Lineage() {
		dbName := cliContext.String("database")
		if dbName != "" && dbName != "<database_name>" {
			return dbName
		}
	}

	var dbName string

	if os.Getenv("SCALINGO_DATABASE") != "" {
		dbName = os.Getenv("SCALINGO_DATABASE")
	}

	return dbName
}

// extractAppNameFromGitRemote parses a Git remote and return the app name extracted
// out of it. The Git remote URL may look like:
// - SSH on a custom port: ssh://git@host:port/appName.git
// - GitHub: git@github.com:owner/appName.git
func extractAppNameFromGitRemote(url string) string {
	splittedURL := strings.Split(strings.TrimSuffix(url, ".git"), ":")
	appName := splittedURL[len(splittedURL)-1]
	// appName may contain "port/appName" or "owner/appName". We keep the part
	// after the slash.
	i := strings.LastIndex(appName, "/")
	if i != -1 {
		appName = appName[i+1:]
	}

	return appName
}
