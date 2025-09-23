// Package cmd gathers the configuration of all commands (names, flags, default
// values etc.) of the Scalingo CLI
package cmd

import (
	"bytes"
	"context"
	"os"
	"text/template"

	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/session"
	"github.com/Scalingo/go-scalingo/v8"
	"github.com/Scalingo/go-scalingo/v8/debug"
)

type AppCommands struct {
	commands []*cli.Command
}

type Command struct {
	*cli.Command
	// Regional flag not available if Global is true
	Global bool
}

type CommandDescription struct {
	// Mandatory description of the command
	Description string

	// Examples on how to use the command
	Examples []string

	//
	SeeAlso []string
}

func (command CommandDescription) Render() string {
	buf := &bytes.Buffer{}
	commandDescriptionTemplate := `{{.Description}}{{ if .Examples }}

Example{{with $length := len .Examples}}{{if ne 1 $length}}s{{end}}{{end}}{{ range .Examples }}
  $ {{ . }}{{ end }}{{ end }}{{ if .SeeAlso }}

# See also{{ range .SeeAlso }} '{{ . }}'{{ end }}{{ end }}`

	template, err := template.New("documentation").Parse(commandDescriptionTemplate)

	if err != nil {
		return ""
	}

	if err := template.Execute(buf, command); err != nil {
		return ""
	}

	return buf.String()
}

func (cmds *AppCommands) addCommand(cmd Command) {
	cmds.commands = append(cmds.commands, cmd.Command)

	// This argument disables the global help command, but doesn't disable the flag.
	cmd.Command.HideHelpCommand = true

	// Global commands are simply added to the list of commands
	if cmd.Global {
		return
	}

	// Regional commands are modified before being added to the list of commands.
	regionFlag := &cli.StringFlag{Name: "region", Value: "", Usage: "Name of the region to use"}
	cmd.Command.Flags = append(cmd.Command.Flags, regionFlag)
	action := cmd.Command.Action
	cmd.Command.Action = regionalCommandAction(action)
}

func regionalCommandAction(action cli.ActionFunc) cli.ActionFunc {
	return func(ctx context.Context, c *cli.Command) error {
		token := os.Getenv("SCALINGO_API_TOKEN")

		currentUser, err := config.C.CurrentUser(ctx)
		if err != nil || currentUser == nil {
			err := session.Login(ctx, session.LoginOpts{APIToken: token})
			if err != nil {
				errorQuit(ctx, err)
			}
		}

		regions, err := config.EnsureRegionsCache(ctx, config.C, config.GetRegionOpts{
			Token: token,
		})
		if err != nil {
			errorQuit(ctx, err)
		}
		currentRegion := regionNameFromFlags(c)

		// Detecting Region from git remote
		if currentRegion == "" {
			currentRegion = detect.GetRegionFromGitRemote(c, &regions)
		}

		if config.C.ScalingoRegion == "" && currentRegion == "" {
			region := getDefaultRegion(regions)
			debug.Printf("[Regions] Use the default region '%s'\n", region.Name)
			currentRegion = region.Name
		}

		if currentRegion != "" {
			config.C.ScalingoRegion = currentRegion
		}

		return action(ctx, c)
	}
}

func getDefaultRegion(regionsCache config.RegionsCache) scalingo.Region {
	defaultRegion := regionsCache.Regions[0]
	for _, region := range regionsCache.Regions {
		if region.Default {
			defaultRegion = region
			break
		}
	}
	return defaultRegion
}

func (cmds *AppCommands) Commands() []*cli.Command {
	return cmds.commands
}

func NewAppCommands() *AppCommands {
	cmds := AppCommands{}
	for _, cmd := range regionalCommands {
		cmds.addCommand(Command{Command: cmd})
	}
	for _, cmd := range globalCommands {
		cmds.addCommand(Command{Global: true, Command: cmd})
	}
	return &cmds
}

var (
	regionalCommands = []*cli.Command{
		// Apps
		&appsCommand,
		&CreateCommand,
		&DestroyCommand,
		&renameCommand,
		&appsInfoCommand,
		&openCommand,
		&dashboardCommand,
		&appsProjectSetCommand,

		// Apps Actions
		&logsCommand,
		&logsArchivesCommand,
		&runCommand,
		&oneOffStopCommand,

		// Apps Process Actions
		&psCommand,
		&scaleCommand,
		&RestartCommand,
		&sendSignalCommand,

		// Routing Settings
		&forceHTTPSCommand,
		&stickySessionCommand,
		&routerLogsCommand,
		&setCanonicalDomainCommand,
		&unsetCanonicalDomainCommand,

		// Events
		&UserTimelineCommand,
		&TimelineCommand,

		// Environment
		&envCommand,
		&envGetCommand,
		&envSetCommand,
		&envUnsetCommand,

		// Domains
		&DomainsListCommand,
		&DomainsAddCommand,
		&DomainsRemoveCommand,
		&DomainsSSLCommand,

		// Deployments
		&deploymentsListCommand,
		&deploymentLogCommand,
		&deploymentFollowCommand,
		&deploymentDeployCommand,
		&deploymentCacheResetCommand,

		// Collaborators
		&CollaboratorsListCommand,
		&CollaboratorsAddCommand,
		&CollaboratorsRemoveCommand,
		&CollaboratorsUpdateCommand,

		// Stacks
		&stacksListCommand,
		&stacksSetCommand,

		// Addons
		&AddonProvidersListCommand,
		&AddonProvidersPlansCommand,
		&addonsListCommand,
		&addonsAddCommand,
		&addonsRemoveCommand,
		&addonsUpgradeCommand,
		&addonsInfoCommand,
		&addonsConfigCommand,

		// Integration Link
		&integrationLinkShowCommand,
		&integrationLinkCreateCommand,
		&integrationLinkUpdateCommand,
		&integrationLinkDeleteCommand,
		&integrationLinkManualDeployCommand,
		&integrationLinkManualReviewAppCommand,

		// Review Apps
		&reviewAppsShowCommand,

		// Notifiers
		&NotifiersListCommand,
		&NotifiersDetailsCommand,
		&NotifiersAddCommand,
		&NotifiersUpdateCommand,
		&NotifiersRemoveCommand,

		// Notification platforms
		&NotificationPlatformListCommand,

		// DB Access
		&DbTunnelCommand,
		&RedisConsoleCommand,
		&MongoConsoleCommand,
		&MySQLConsoleCommand,
		&PgSQLConsoleCommand,
		&InfluxDBConsoleCommand,

		// Databases
		&databaseBackupsConfig,
		&databaseEnableFeature,
		&databaseDisableFeature,
		&databaseListUsers,
		&databaseDeleteUser,
		&databaseCreateUser,
		&databaseUpdateUserPassword,

		// Maintenance
		&databaseMaintenanceList,
		&databaseMaintenanceInfo,

		// Backups
		&backupsListCommand,
		&backupsCreateCommand,
		&backupsDownloadCommand,
		&backupDownloadCommand,

		// Alerts
		&alertsListCommand,
		&alertsAddCommand,
		&alertsUpdateCommand,
		&alertsEnableCommand,
		&alertsDisableCommand,
		&alertsRemoveCommand,

		// Stats
		&StatsCommand,

		// Autoscalers
		&autoscalersListCommand,
		&autoscalersAddCommand,
		&autoscalersRemoveCommand,
		&autoscalersUpdateCommand,
		&autoscalersDisableCommand,
		&autoscalersEnableCommand,

		// Log drains
		&logDrainsAddCommand,
		&logDrainsListCommand,
		&logDrainsRemoveCommand,

		&gitSetup,
		&gitShow,

		// Cron tasks
		&cronTasksListCommand,

		// Projects
		&projectsListCommand,
		&projectsAddCommand,
		&projectsUpdateCommand,
		&projectsRemoveCommand,
	}

	globalCommands = []*cli.Command{
		// SSH keys
		&listSSHKeyCommand,
		&addSSHKeyCommand,
		&removeSSHKeyCommand,

		&integrationsListCommand,
		&integrationsAddCommand,
		&integrationsDeleteCommand,
		&integrationsImportKeysCommand,

		// Sessions
		&LoginCommand,
		&LogoutCommand,
		&RegionsListCommand,
		&ConfigCommand,
		&selfCommand,

		// Version
		&UpdateCommand,

		// Changelog
		&changelogCommand,

		// Help
		&HelpCommand,
	}
)
