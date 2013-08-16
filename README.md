Appsdeck-CLI
============

Command line utility to manage its appsdeck application.

```
NAME:
   Appsdeck Client - Manage your apps and containers

USAGE:
   Appsdeck Client [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
   logs, l	Print logs of current app
   run, r	Run any command for your app
   apps, a	Manage your apps
   logout	Logout from Appsdeck
   help, h	Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --app '<name>'	Name of the app
   --version		print the version
   --help, -h		show help
```

Dev usage
---------

Define (example) :

* `APPSDECK_LOG=127.0.0.1:10004`
* `APPSDECK_API=127.0.0.1:3000`
