Scalingo-CLI v1.8.0
===================

This repository contains the command line utility for the public PaaS Scalingo

https://scalingo.com

## How to build?

The project is using go, then you need a running go environment: [Official documentation](https://golang.org/doc/install)

Once that's done, all you have to do is to 'go get' the project, with the following command:

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
http_proxy=http://<proxy host>:<proxy port>
https_proxy=https://<proxy host>:<proxy port>
```

## Command help

```

NAME:
   Scalingo Client - Manage your apps and containers

USAGE:
   scalingo [global options] command [command options] [arguments...]

VERSION:
   1.7.0

AUTHOR:
   Scalingo Team <hello@scalingo.com>

COMMANDS:
     help  Shows a list of commands or help for one command

   Addons:
     addons          List used add-ons
     addons-add      Provision an add-on for your application
     addons-remove   Remove an existing addon from your app
     addons-upgrade  Upgrade or downgrade an add-on attached to your app

   Addons - Global:
     addons-list   List all addons
     addons-plans  List plans

   App Management:
     destroy            Destroy an app /!\
     rename             Rename an application
     logs, l            Get the logs of your applications
     logs-archives, la  Get the logs archives of your applications
     run, r             Run any command for your app
     ps                 Display your application running processes
     scale, s           Scale your application instantly
     restart            Restart processes of your app
     db-tunnel          Create an encrypted connection to access your database

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
     deployments        List app deployments
     deployment-logs    View deployment logs
     deployment-follow  Follow deployment event stream
     deploy             Trigger a deployment by archive

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
     signup     Create your Scalingo account

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

GLOBAL OPTIONS:
   --app value, -a value     Name of the app (default: "<name>") [$SCALINGO_APP]
   --remote value, -r value  Name of the remote (default: "scalingo")
   --version, -v             print the version
```
