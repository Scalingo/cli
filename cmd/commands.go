package cmd

import (
	"github.com/urfave/cli"
)

var (
	Commands = []cli.Command{
		// Apps
		AppsCommand,
		CreateCommand,
		DestroyCommand,
		RenameCommand,

		// Apps Actions
		LogsCommand,
		LogsArchivesCommand,
		RunCommand,

		// Apps Process Actions
		PsCommand,
		ScaleCommand,
		RestartCommand,

		// Events
		UserTimelineCommand,
		TimelineCommand,

		// Environment
		EnvCommand,
		EnvSetCommand,
		EnvUnsetCommand,

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

		// Collaborators
		CollaboratorsListCommand,
		CollaboratorsAddCommand,
		CollaboratorsRemoveCommand,

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

		// TODO: Alerts
		alertsListCommand,
		alertsAddCommand,

		// Stats
		StatsCommand,

		// SSH keys
		ListSSHKeyCommand,
		AddSSHKeyCommand,
		RemoveSSHKeyCommand,

		// Sessions
		LoginCommand,
		LogoutCommand,
		SignUpCommand,

		// Version
		UpdateCommand,

		// Help
		HelpCommand,
	}
)
