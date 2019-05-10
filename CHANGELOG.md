# Changelog

### To be Released

* [alerts] Add support for the duration_before_trigger attribute [#407](https://github.com/Scalingo/cli/pull/407)
* [domains] Handle Let's Encrypt certificate status [#410](https://github.com/Scalingo/cli/pull/410)

### 1.10.1

* Wrong default URL for the database API [#403](https://github.com/Scalingo/cli/pull/403)

### 1.10.0

* [db] Add database logs #398
* [db] Add backups and backup-download commands #397
* Add force HTTPS, sticky session and canonical domain #344
* [commands] Update scale commands to accept flags to create an autoscaler #339
* [env-set] Advise to restart after setting the environment #373
* Add `--password-only` flag to `scalingo login`: to change account when a SSH key is defined #351
* [self] Add `self` (and `whoami` alias): to know which user is connected in case of multi-account #350
* Better error message if 401 Unauthorized #352
* Fix English wording #385
* Fix missing spaces #388
* Fix new format of SSH key 'invalid type' error #389
* [db-tunnel] Only display parsable information on stdout, the rest on stderr #396

### 1.9.0

* [alerts] Add all command to CRUD alerts #346

```
   Alerts:
     alerts          List the alerts of an application
     alerts-add      Add an alert to an application
     alerts-update   Update an alert
     alerts-enable   Enable an alert
     alerts-disable  Disable an alert
     alerts-remove   Remove an alert from an application
```

* [notifiers] Add the ability to configure email notifiers with custom emails and collaborators #366

```
$ scalingo -a my-app notifiers-add --platform email --send-all-events --name email-notif-1 --email notifications@example.com --collaborator username1 --collaborator username2
+-----------------+------------------------------------------+
| ID              | no-bd6ea457-ccf9-45c5-8767-04225fdb1018  |
| Type            | email                                    |
| Name            | email-notif-1                            |
| Enabled         | true                                     |
| Send all events | true                                     |
| Emails          | [notifications@example.com]              |
| User_ids        | [us-38246321-111f-4b54-a22c-12b04548c55f |
|                 | us-0784238c-b422-4a79-8760-f3ffbd10705c] |
+-----------------+------------------------------------------+
```

* [deployment] Add command to reset deployment cache #358

```
$ scalingo -a my-app deployment-delete-cache
-----> Deployment cache successfully deleted
```

* [deployment] Fix deploy/deployments-follow log streaming when multiple deployments are running #359
* [update] Add the ability to disable the update checker with the environment variable `DISABLE_UPDATE_CHECKER=true` #361
* [global] Correctly display help when command syntax is not respected #367
* [logs] Bugfix: consider Ctrl^C the default way to stop `logs -f` command, it's not an error #368

### 1.8.0

* [logs] accepts filter for postdeploy and one-off container
* [authentication] use of the new authentication API `auth.scalingo.com`
* [help] fix missing help with unknown commsnds #365

### 1.7.0

* [Commands] Add `rename` command to rename an application
* [One-off] Better inactivity timeout error message
* [DB Console] Add support for TLS connection to databases
* [Bugfix] Bad autocompletion on -a, --app, -r, --remote flags when they are the first argument of a command
* [Bugfix] TTY size was not sent when launching a `run` command


### v1.6.0

* [Mongo Console] Add replicaset support to correctly connect to them #306
* [Notifiers] Add Notifiers related commands #303 #301 #297 #296:

```
     notifiers          List your notifiers
     notifiers-details  Show details of your notifiers
     notifiers-add      Add a notifier for your application
     notifiers-update   Update a notifier
     notifiers-remove   Remove an existing notifier from your app
```

* [Notifications] Feature removed, replaced by notifiers, all the notifications have been migrated to notifiers #307
* [Internals] Migrate to original `urfave/cli` instead of using our own fork of the library #290
* [Update] Add timeout in update checking to avoid the CLI to freeze when GitHub is down for instance #274
* [Auth] When authentication file is corrupted, recreate a new one instead of crashing #283
* [Logs-archive] Logs archives are now listable and downloadable from the CLI #289
* [Logs] Lines are now colored according to the source of the line #286

### v1.5.1

* [Feature] Authenticate your request using the environment by using the environment variable `SCALINGO_API_TOKEN` #291

### v1.5.0

* [Feature] Add `deploy` command to deploy a tarball or a war archive directly

```
scalingo deploy archive.tar.gz
scalingo deploy project.war
scalingo deploy https://github.com/Scalingo/sample-go-martini/archive/master.tar.gz
```

### v1.4.1

* [Fix] Fix error message when a user tries to break its free trial before the end #458
* [Feature] Add `influxdb-console` to run an influxdb interactive shell in a one-off container

### v1.4.0 - 16/11/2016

* [Feature] Add timeline and user-timeline to display per are of user-global activities #235
* [Feature] Add list, remove and add commands for notifications
* [Feature] Add `deployments` command to get the a deployments list for an application #222 #234
* [Feature] Add `deployment-logs` command to get logs for a specific deployment
* [Feature] Add `deployment-follow` command to follow the deployment stream for an application
* [Feature - Login] Automatically try SSH with ssh-agent if available #262
* [Feature - Create] --buildpack flag to specify a custom buildpack
* [Fix] Fix error handling when an addon fails to get provisioned #252
* [Fix] Fix error display when an application doesn't have any log available #249
* [Fix] Fix error display when connection to the SSH server fails #242
* [Fix - Windows] Password typing error on windows (ReadConsoleInput error) #237
* [Fix] Login command logs twice #258
* [Fix - MacOS Sierra] Build with go 1.7, fulle compatible with Sierra

### v1.3.1 - 02/05/2016

* [Bugfix - Auth] Fix authentication configuration for --ssh or --apikey, two attempts were necessary #208 #209

### v1.3.0 - 01/04/2016

* [Feature - Auth] Authentication with API key or SSH key (--ssh or --api-key flags) #196 #200
* [Feature - Auth] New format of configuration file for authentication, auto migration. #200
* [Feature - Scale] Possibility to scale with relative operator (i.e. web:+1) #197 #198
* [Feature - Run] --type to directly run a command defined by a Procfile line #185 #207
* [Feature - Run] --silent flag to remove any noise and only get the one-off command output #191
* [Enhancement - Run] Display output on stderr to be able to drop it to /dev/null #190
* [Enhancement - Run] Exit code of one-off container is now forward as exit code of the CLI #203 #205
* [Bugfix - Stats] Fix computation of percentage for higher bound value
* [Bugfix - Run] Accept pipes and redirections as input for one-off containers #199 #206
* [Bugfix - Env] Remove arguments validation, that's server role, and it changes sometimes
* [Bugfix - Env] Add quotes in output of env-set to avoid copy/paste problem with the final period
* [Bugfix - Scale] Fix error management when application is already restarting or scaling #195
* [Bugfix - Tunnel] Fix panic when authentication fails when building SSH tunnel
* [Bugfix - Tunnel] Fix double error handling when binding local port #202
* [Bugfix] Fix install script on Mac OS X El Capitan 10.11

Contributors
- leo@scalingo.com <Leo Unbekandt>
- mail2tevin@gmail.com <Tevin Zhang>

### v1.2.0 - 20/11/2015

* [Feature - DB Tunnel] Reconnect automatically in case of connection problem
* [Feature - DB Tunnel] Default port at 10000, if not available 10001 etc.
* [Feature - One-off] More verbose output and spinner when starting a one-off container #180 #184
* [Feature - Logs] Automatically reconnect to logs streaming if anything wrong happen #182
* [Feature] Add `stats` command to get containers CPU and memory metrics
* [Bugfix] Fix delete command (app name wasn't read correctly) #177

### v1.1.3 - 29/10/2015

* [Bugfix] Authentication problem when auth file doesn't exist

### v1.1.2 - 23/10/2015

* [Feature] Show suggestions to wrong commands #164
* [Feature] Add `DISABLE_INTERACTIVE` environment variable to disable blocking user input #146
* [Feature - Completion] Enable completion on restart command #158 #159
* [Bugfix] Login on Windows 10 when using git bash #171 #160
* [Bugfix] Fix error when upgrading addon #168 #170
* [Bugfix] User friendly login prompt in case of "No account" #152
* [Bugfix] Destroy command requesting API to know if app exists or not before asking for confirmation #161 #162 #155 #151
* [Bugfix] Do not display wrong completion when user is not logged in #146 #142
* [Refactoring] Extract Scalingo API functions to an external package github.com/Scalingo/go-scalingo #150
* [Refactoring] Use API endpoint to update multiple environment variables at once #153

### v1.1.1 - 22/09/2015

* [Feature] Build in Linux ARM #145
* [Feature - Completion] Add local cache of applications when using completion on them, avoid heavy unrequired API requests #141
* [Feature - Completion] Completion of the `--remote` flag #139
* [Optimisation - Completion] Completion of `collaborators-add` command is now quicker (×2 - ×4) #137
* [Bugfix - Completion] Do not display error in autocompletion when unlogged #142
* [Bugfix] Fix regression, small flags were not working anymore #144 #147

### v1.1.0 - 17/09/2015

* [Feature - CLI] Setup Bash and ZSH completion thanks to codegangsta/cli helpers #127
* [Feature - CLI] Add -r/--remote flag to specify a GIT remote instead of an app name #89
* [Feature - CLI] Add -r/--remote flag to the `create` subcommand to specify an alternative git remote name (default is `scalingo`) #129
* [Feature - Log] Add -F/--filter flag to filter log output by container types #118
* [Bugfix - Run] Fix parsing of environment variables (flag -e) #119
* [Bugfix - Mongo Console] Do not try to connect to the oplog user anymore (when enabled) #117
* [Bugfix - Logs] Stream is cut with an 'invalid JSON' error, fixed by increasing the buffer size #135
* [Bugfix - Tunnel] Error when the connection to the database failed, a panic could happen

### v1.0.0 - 06/05/2015

* [Feature - Databases] Add helper to run interactive console for MySQL, PostgreSQL, MongoDB and Redis #111
* [Feature - Collaborators] Handle collaborators directly from the command line client #113
* [Feature - Proxy] Add support and documentation about how to use a HTTPS proxy #104 #110
* [Refactoring - API calls] Completely refactor error management for Scalingo API calls #112
* [Improvement - SSL] Embed Scalingo new SSL certificate SHA-256 only #109
* [Bugfix - Macos] #105 #114
* [Bugfix - Logs] No more weird error message when no log is available for an app #108
* [Bugfix - Logs] Use of websocket for log streaming #86 #115 #116
* [Bugfix - Windows] Babun shell compatibility #106

### v1.0.0-rc1 - 15/04/2015

* [Feature] Modify size of containers with `scalingo scale` - #102
* [Bugfix] Fix ssh-agent error when no private key is available - Fixed #100
* [Bugfix] Fix domain-add issue. (error about domain.crt file) - Fixed #98
* [Bugfix] Fix addon plans description, no more HTML in them  - #96
* [Bugfix] Correctly handle db-tunnel when alias is given as argument - Fixed #93

### v1.0.0-beta1 - 08/03/2015

* Windows, password: don't display password in clear
* Windows, db-tunnel: Correctly handle SSH key path with -i flag
* Send OS to one-off containers (for prompt handling, useful for Windows)
* Fix EOF error when writing the password
* Fix authentication request to adapt the API
* Correctly handle 402 errors (payment method required) #90
* Project is go gettable `go get github.com/Scalingo/cli`
* Fix GIT remote detection #89
* Correctly handle 404 Error, display clearer messages #91
* More documentation for the `run` command - Fixed #79
* Rewrite API client package, remove unsafe code - Fixed #80
* Allow environment variable name or value for `db-tunnel` as argument
* Extended help for `db-tunnel` - Fixed #85
* Ctrl^C doesn't kill an `run` command anymore - Fixed #83
* --app flag can be written everywhere in the command line - Fixed #10
* Use SSH agent if possible to get SSH credentials
* Correcty handle encrypted SSH keys (AES-CBC and DES-ECE2) - Fixed #76, #77

### v1.0.0-alpha4 - 22/01/2014

* Rename `Processes` to `Containers` to fit the API
* `new-command`: `login` command
* `logout`: Clean credentials deletion for multiple APIs (dev feature)
* `logs`: Do not encode HTML entities (logs server don't escape html entities anymore)
* `addons-*`: Adapt to new addons API endpoints
* `domains-*`: Adapt to new domains API endpoints
* `db-tunnel`: Handle encrypted SSH key

### v1.0.0-alpha3 - 21/12/2014

* Fix credential storage issue #72 #73
* Fix wrong help for command db-tunnel #74
* Fix logfile open operation on MacOS #70
* Build Windows version on Windows with CGO #71
* Build Mac OS verison on Mac OS with CGO #71

### v1.0.0-alpha2 - 15/12/2014

* Move addons subcommands at top level for better visibility

### v1.0.0-alpha1 - 14/12/2014

* Initial release
* Opensourcing of the project
