Scalingo-CLI v1.1.1
=======================

This repository contains the command line utility for the public PaaS Scalingo

https://scalingo.com

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
   Scalingo Client [global options] command [command options] [arguments...]

VERSION:
   1.1.0

AUTHOR:
  Scalingo Team - <hello@scalingo.com>

COMMANDS:
  Addons:
    addons              List used add-ons
    addons-add          Provision an add-on for your application
    addons-remove       Remove an existing addon from your app
    addons-upgrade      Upgrade or downgrade an add-on attached to your app

  Addons - Global:
    addons-list         List all addons
    addons-plans        List plans

  App Management:
    logs, l     Get the logs of your applications
    run, r      Run any command for your app
    ps          Display your application running processes
    scale, s    Scale your application instantly
    restart     Restart processes of your app
    db-tunnel   Create an encrypted connection to access your database

  CLI Internals:
    version     Display current version
    update      Update 'scalingo' client

  Collaborators:
    collaborators               List the collaborators of an application
    collaborators-add           Invite someone to work on an application
    collaborators-remove        Revoke permission to collaborate on an application

  Custom Domains:
    domains             List the domains of an application
    domains-add         Add a custom domain to an application
    domains-remove      Remove a custom domain from an application
    domains-ssl         Enable or disable SSL for your custom domains

  Databases:
    redis-console       Run an interactive console with your redis addon
    mongo-console       Run an interactive console with your MongoDB addon
    mysql-console       Run an interactive console with your MySQL addon
    pgsql-console       Run an interactive console with your PostgreSQL addon

  Environment:
    env         Display the environment of your apps
    env-set     Set the environment variables of your apps
    env-unset   Unset environment variables of your apps

  Global:
    apps, a     List your apps
    create, c   Create a new app
    destroy, d  Destroy an app /!\
    login       Login to Scalingo platform
    logout      Logout from Scalingo
    signup      Create your Scalingo account

  Public SSH Keys:
    keys        List your SSH public keys
    keys-add    Add a public SSH key to deploy your apps
    keys-remove Remove a public SSH key

GLOBAL OPTIONS:
   --app, -a '<name>'   Name of the app [$SCALINGO_APP]
   --help, -h           show help
   --version, -v        print the version
```
