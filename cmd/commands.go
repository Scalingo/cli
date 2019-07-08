package cmd

import (
	"github.com/Scalingo/cli/config"
	scalingo "github.com/Scalingo/go-scalingo"
	"github.com/Scalingo/go-scalingo/debug"
	"github.com/urfave/cli"
)

type AppCommands struct {
	commands []cli.Command
}

type Command struct {
	cli.Command
	// Regional flag not available if Global is true
	Global bool
}

func (cmds *AppCommands) AddCommand(cmd Command) {
	if !cmd.Global {
		regionFlag := cli.StringFlag{Name: "region", Value: "", Usage: "Name of the region to use"}
		cmd.Command.Flags = append(cmd.Command.Flags, regionFlag)
	}
	action := cmd.Command.Action.(func(c *cli.Context))
	cmd.Command.Action = func(c *cli.Context) {
		regions, err := config.EnsureRegionsCache(config.C, config.GetRegionOpts{})
		if err != nil {
			panic(err)
		}
		currentRegion := c.GlobalString("region")
		if currentRegion == "" {
			currentRegion = c.String("region")
		}

		if config.C.ScalingoRegion == "" && currentRegion == "" {
			region := getDefaultRegion(regions)
			debug.Printf("[Regions] Use the default region '%s'\n", region.Name)
			currentRegion = region.Name
		}

		if currentRegion != "" {
			config.C.ScalingoRegion = currentRegion
		}
		action(c)
	}
	cmds.commands = append(cmds.commands, cmd.Command)
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

func (cmds *AppCommands) Commands() []cli.Command {
	return cmds.commands
}

func NewAppCommands() *AppCommands {
	cmds := AppCommands{}
	for _, cmd := range regionalCommands {
		cmds.AddCommand(Command{Command: cmd})
	}
	for _, cmd := range globalCommands {
		cmds.AddCommand(Command{Global: true, Command: cmd})
	}
	return &cmds
}

var (
	regionalCommands = []cli.Command{
		// Apps
		appsCommand,
		CreateCommand,
		DestroyCommand,
		renameCommand,
		appsInfoCommand,

		// Apps Actions
		LogsCommand,
		LogsArchivesCommand,
		RunCommand,

		// Apps Process Actions
		psCommand,
		scaleCommand,
		RestartCommand,

		// Routing Settings
		forceHTTPSCommand,
		stickySessionCommand,
		setCanonicalDomainCommand,
		unsetCanonicalDomainCommand,

		// Events
		UserTimelineCommand,
		TimelineCommand,

		// Environment
		envCommand,
		envSetCommand,
		envUnsetCommand,

		// Domains
		DomainsListCommand,
		DomainsAddCommand,
		DomainsRemoveCommand,
		DomainsSSLCommand,

		// Deployments
		DeploymentsListCommand,
		DeploymentLogCommand,
		DeploymentFollowCommand,
		DeploymentDeployCommand,
		DeploymentCacheResetCommand,

		// Collaborators
		CollaboratorsListCommand,
		CollaboratorsAddCommand,
		CollaboratorsRemoveCommand,

		// Stacks
		stacksListCommand,
		stacksSetCommand,

		// Addons
		AddonProvidersListCommand,
		AddonProvidersPlansCommand,
		AddonsListCommand,
		AddonsAddCommand,
		AddonsRemoveCommand,
		AddonsUpgradeCommand,

		// Notifiers
		NotifiersListCommand,
		NotifiersDetailsCommand,
		NotifiersAddCommand,
		NotifiersUpdateCommand,
		NotifiersRemoveCommand,

		// Notification platforms
		NotificationPlatformListCommand,

		// DB Access
		DbTunnelCommand,
		RedisConsoleCommand,
		MongoConsoleCommand,
		MySQLConsoleCommand,
		PgSQLConsoleCommand,
		InfluxDBConsoleCommand,

		// Backups
		BackupListCommand,
		BackupDownloadCommand,

		// Alerts
		alertsListCommand,
		alertsAddCommand,
		alertsUpdateCommand,
		alertsEnableCommand,
		alertsDisableCommand,
		alertsRemoveCommand,

		// Stats
		StatsCommand,

		// Autoscalers
		autoscalersListCommand,
		autoscalersAddCommand,
		autoscalersRemoveCommand,
		autoscalersUpdateCommand,
		autoscalersDisableCommand,
		autoscalersEnableCommand,

		gitSetup,
		gitShow,
	}

	globalCommands = []cli.Command{
		// SSH keys
		ListSSHKeyCommand,
		AddSSHKeyCommand,
		RemoveSSHKeyCommand,

		// Sessions
		LoginCommand,
		LogoutCommand,
		RegionsListCommand,
		ConfigCommand,
		selfCommand,
		whoamiCommand,

		// Version
		UpdateCommand,

		// Help
		HelpCommand,
	}
)
