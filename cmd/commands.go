package cmd

import (
	"github.com/urfave/cli"
)

var (
	Commands = []cli.Command{
		// Apps
		appsCommand,
		CreateCommand,
		DestroyCommand,
		RenameCommand,

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

		// TODO: Alerts
		alertsListCommand,
		alertsAddCommand,
		alertsUpdateCommand,
		alertsEnableCommand,
		alertsDisableCommand,
		alertsRemoveCommand,

		// Stats
		StatsCommand,

		// SSH keys
		ListSSHKeyCommand,
		AddSSHKeyCommand,
		RemoveSSHKeyCommand,

		// Autoscalers
		autoscalersListCommand,
		autoscalersAddCommand,
		autoscalersRemoveCommand,
		autoscalersUpdateCommand,
		autoscalersDisableCommand,
		autoscalersEnableCommand,

		// Sessions
		LoginCommand,
		LogoutCommand,
		RegionsListCommand,
		RegionsSetCommand,
		selfCommand,
		whoamiCommand, // `self` alias

		// Version
		UpdateCommand,

		// Help
		HelpCommand,
	}
)
