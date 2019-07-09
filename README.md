# Scalingo-CLI v1.14.0

This repository contains the command line utility for the public PaaS Scalingo

https://scalingo.com

## How to build?

The project is using Go, then you need a running Go environment: [Official documentation](https://golang.org/doc/install)

Once that's done, all you have to do is to `go get` the project, with the following command:

```
go get github.com/Scalingo/cli/scalingo
```

That's it, you've build the latest version of the Scalingo CLI (the binary will be present in `$GOPATH/bin/scalingo`)


## How to upgrade?

```
go get -u github.com/Scalingo/cli/scalingo
```

## Run behind a proxy

You have to setup the following environment variables:

```
HTTP_PROXY=http://<proxy host>:<proxy port>
HTTPS_PROXY=https://<proxy host>:<proxy port>
```

## Disable update checking

By default the CLI is making an HTTP request to learn if a newer version is available.
To disable this feature, define the environment variable:

```
DISABLE_UPDATE_CHECKER=true
```

## Command help

```

NAME:
   Scalingo Client - Manage your apps and containers

USAGE:
   scalingo [global options] command [command options] [arguments...]

VERSION:
   1.14.0

AUTHOR:
   Scalingo Team <hello@scalingo.com>

COMMANDS:
     help  Shows a list of commands or help for one command

   Addons:
     addons           List used add-ons
     addons-add       Provision an add-on for your application
     addons-remove    Remove an existing addon from your app
     addons-upgrade   Upgrade or downgrade an add-on attached to your app
     backups          List backups for an addon
     backup-download  Download a backup

   Addons - Global:
     addons-list   List all addons
     addons-plans  List plans

   Alerts:
     alerts          List the alerts of an application
     alerts-add      Add an alert to an application
     alerts-update   Update an alert
     alerts-enable   Enable an alert
     alerts-disable  Disable an alert
     alerts-remove   Remove an alert from an application

   App Management:
     destroy                 Destroy an app /!\
     rename                  Rename an application
     logs, l                 Get the logs of your applications
     logs-archives, la       Get the logs archives of your applications
     run, r                  Run any command for your app
     ps                      Display your application running processes
     scale, s                Scale your application instantly
     restart                 Restart processes of your app
     force-https
     sticky-session
     set-canonical-domain    Set a canonical domain.
     unset-canonical-domain  Unset a canonical domain.
     db-tunnel               Create an encrypted connection to access your database

   Autoscalers:
     autoscalers          List the autoscalers of an application
     autoscalers-add      Add an autoscaler to an application
     autoscalers-remove   Remove an autoscaler from an application
     autoscalers-update   Update an autoscaler
     autoscalers-disable  Disable an autoscaler
     autoscalers-enable   Enable an autoscaler

   CLI Internals:
     update  Update 'scalingo' client

   Collaborators:
     collaborators         List the collaborators of an application
     collaborators-add     Invite someone to work on an application
     collaborators-remove  Revoke permission to collaborate on an application

   Custom Domains:
     domains         List the domains of an application
     domains-add     Add a custom domain to an application
     domains-remove  Remove a custom domain from an application
     domains-ssl     Enable or disable SSL for your custom domains

   Databases:
     redis-console     Run an interactive console with your Redis addon
     mongo-console     Run an interactive console with your MongoDB addon
     mysql-console     Run an interactive console with your MySQL addon
     pgsql-console     Run an interactive console with your PostgreSQL addon
     influxdb-console  Run an interactive console with your InfluxDB addon

   Deployment:
     deployments              List app deployments
     deployment-logs          View deployment logs
     deployment-follow        Follow deployment event stream
     deploy                   Trigger a deployment by archive
     deployment-delete-cache  Reset deployment cache

   Display metrics of the running containers:
     stats  Display metrics of the currently running containers

   Environment:
     env        Display the environment of your apps
     env-set    Set the environment variables of your apps
     env-unset  Unset environment variables of your apps

   Events:
     user-timeline  List the events you have done on the platform
     timeline       List the actions related to a given app

   Global:
     apps       List your apps
     create, c  Create a new app
     login      Login to Scalingo platform
     logout     Logout from Scalingo
     regions    List available regions
     config     Configure the CLI
     signup     Create your Scalingo account
     self       Get the logged in profile
     whoami     Get the logged in profile

   Notifiers:
     notifiers          List your notifiers
     notifiers-details  Show details of your notifiers
     notifiers-add      Add a notifier for your application
     notifiers-update   Update a notifier
     notifiers-remove   Remove an existing notifier from your app

   Notifiers - Global:
     notification-platforms  List all notification platforms

   Public SSH Keys:
     keys         List your SSH public keys
     keys-add     Add a public SSH key to deploy your apps
     keys-remove  Remove a public SSH key

   Runtime Stacks:
     stacks      List the available runtime stacks
     stacks-set  Set the runtime stack of an app

GLOBAL OPTIONS:
   --addon value             ID of the current addon (default: "<addon_id>") [$SCALINGO_ADDON]
   --app value, -a value     Name of the app (default: "<name>") [$SCALINGO_APP]
   --remote value, -r value  Name of the remote (default: "scalingo")
   --region value            Name of the region to use
   --version, -v             print the version
```

## Development setup

In order to build the current development version and use it against the development services:

```bash
cd scalingo
go build .
SCALINGO_API_URL=http://172.17.0.1:3001 SCALINGO_AUTH_URL=http://172.17.0.1:1234 ./scalingo login --api-token <admin user API token>
```

## Release a new version

Bump new version number in:

- `.goxc.json`
- `CHANGELOG.md`
- `README.md`
- `VERSION`
- `config/version.go`

Build the new version for all platforms with: `./dists/make-release.sh -v 1.14.0`.

Tag and release a new version on GitHub [here](https://github.com/Scalingo/cli/releases/new). Attach
the zip archives created by the `make-release.sh` script to this release.

Last, restart the Scalingo application `cli-download-service`. It serves as cache between GitHub and
our customers for a more efficient check of what is the new CLI version. Type:

```
scalingo -a cli-download-service restart
```

You can now update the [changelog](https://doc.scalingo.com/changelog) and tweet about it!
