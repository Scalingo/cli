Appsdeck-CLI
============

Command line utility to manage its appsdeck application.

```
NAME:
   Appsdeck Client - Manage your apps and containers

USAGE:
   Appsdeck Client [global options] command [command options] [arguments...]

VERSION:
   0.3.4

COMMANDS:
   logs, l	[-n <nblines> | --stream]
   run, r	Run any command for your app
   apps, a	Manage your apps
   logout	Logout from Appsdeck
   create, c	appsdeck create <name>
   destroy, d	appsdeck destroy <id or canonical name>
   help, h	Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --app, -a '<name>'	Name of the app
   --version, -ve	print the version
   --help, -h		show help

```

Dev usage
---------

Define (example) :

* `APPSDECK_API=http://appsdeck.dev`
* `UNSECURE_SSL=true`
* `DEBUG=1`
