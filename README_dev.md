Scalingo-CLI
============

Command line utility to manage its Scalingo apps.

```
NAME:
   Scalingo Client - Manage your apps and containers

USAGE:
   Scalingo Client [global options] command [command options] [arguments...]

VERSION:
   0.5.0

COMMANDS:
   logs, l	[-n <nblines> | --stream]
   run, r	Run any command for your app
   apps, a	Manage your apps
   logout	Logout from Scalingo
   create, c	scalingo create <name>
   destroy, d	scalingo destroy <id or canonical name>
   help, h	Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --app, -a '<name>'	Name of the app
   --version, -ve	print the version
   --help, -h		show help

```

Dev usage
---------

Define (example) :

* `API_URL=http://scalingo.dev`
* `UNSECURE_SSL=true`
* `DEBUG=1`
