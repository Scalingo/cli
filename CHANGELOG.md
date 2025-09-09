# Changelog

## To Be Released

* task(backups-config): update error message when unscheduling periodic backup
* feat(backups): prevent download of OpenSearch and Elasticsearch backups

### 1.38.0

/!\ This change required the [command completion](https://doc.scalingo.com/tools/cli/start#command-completion) to be re-run
* chore(deps): update urfave/cli from v2 to v3

### 1.37.0

* feat(apps): add the project of the applications in `apps` and `apps-info`
* feat(apps): create a new filter `project` on the `apps` command
* feat(apps): add parameter `project-id` on the `create` command
* feat(projects): create all commands to handle `projects`

### 1.36.0

* feat(backups-config): warn user before unscheduling periodic backups
* feat(collaborators): update list to include `is_limited` information
* feat(collaborators): add command to update a collaborator

### 1.35.1

* chore(deps): Bump go-scalingo to v8.1.0

### 1.35.0

* feat(regionmigrations): remove the commands [#1104](https://github.com/Scalingo/cli/pull/1104)
* chore(go) update go version to 1.24 [#1108](https://github.com/Scalingo/cli/pull/1108)
* fix(redis/logs): fix retrieval of logs for redis addon when using `--addon redis` [#1115](https://github.com/Scalingo/cli/pull/1115)

### 1.34.0

* feat(database/users): raise minimum user password length to 24 ([PR#1077](https://github.com/Scalingo/cli/pull/1077))

### 1.33.0

* fix(one-off): remove async field from the run command ([PR#1060](https://github.com/Scalingo/cli/pull/1060))
* feat(database/users): use API DatabaseUserResetPassword method for resetting password

### 1.32.0

* feat(domains): `domains-add` now supports `--no-letsencrypt` flag to completely disable Let's Encrypt certificat generation ([PR#1058](https://github.com/Scalingo/cli/pull/1058))
* fix(cmd): `addons-info` working for non-db addons ([PR#1055](https://github.com/Scalingo/cli/pull/1055))

### 1.31.0

* fix(autocomplete): do not show update message on autocomplete ([PR#1044](https://github.com/Scalingo/cli/pull/1044))
* feat(databases): add databases users update ([PR#1036](https://github.com/Scalingo/cli/pull/1036))

### 1.30.1

* feat(addons): add aliases for database user management (i.e. using format: `database-users-<action>`) ([PR#1031](https://github.com/Scalingo/cli/pull/1031))
* fix(login): regression two-factor token is now correctly asked again ([PR#1037](https://github.com/Scalingo/cli/pull/1037))

### 1.30.0

* feat(addons): add user management commands ([PR#1019](https://github.com/Scalingo/cli/pull/1019))
* chore(term): remove `github.com/andrew-d/go-termutil`, use standard library instead ([PR#974](https://github.com/Scalingo/cli/pull/974))
* feat(cmd): addon can be retrieve from addon type, not only UUID ([PR#983](https://github.com/Scalingo/cli/pull/983))
* feat: add database maintenance info with the new `database-maintenance-info` command ([PR#984](https://github.com/Scalingo/cli/pull/984))
* feat: add database maintenance listing with the new `database-maintenance-list` command ([PR#982](https://github.com/Scalingo/cli/pull/982))
* feat(addons): add maintenance windows manipulation with the new `addon-config` command ([PR#955](https://github.com/Scalingo/cli/pull/955))
* feat(install.sh): verify the archive checksum ([PR#988](https://github.com/Scalingo/cli/pull/988))
* feat(region): more debug logs  ([PR#1007](https://github.com/Scalingo/cli/pull/1007))
* fix(events): `link_scm` data types

### 1.29.1

* revert: refactor: various linter offenses

### 1.29.0

* fix(postgresql): accept database url starting with `postgresql://`
* feat(log-drains): add Logtail
* feat(install.sh): add arm64 to the list of installable architectures ([PR#930](https://github.com/Scalingo/cli/pull/930))
* chore(deps): replace `github.com/ScaleFT/sshkeys` with `golang.org/x/crypto/ssh`
* fix(completion): fix zsh shebang reference
* feat(integration-link): add `--follow` arg to `manual-deploy` command
* fix(cli): return an exit code equal to 1 when a command or flag is malformed
* feat(install.sh): handling word answers as user inputs for install script
* fix(region_migrations): help message for `migration-abort` ([PR#951](https://github.com/Scalingo/cli/pull/951))
* feat(env): set variables from an env file ([PR#950](https://github.com/Scalingo/cli/pull/950))
* feat(install.sh): change default target directory to /opt/homebrew/bin for Apple Silicon

### 1.28.2

* regression(pgsql-console|mongo-console|influx-console|redis-console): Fix regression introduced in 1.28.1, use same async mecanism #913
* refacto(run): Improve code clarity for one-off starting operation and remove annoying extra space #912

### 1.28.1

* feat(one-off): start attached one-off asynchronously

### 1.28.0

* doc: update default stack in help command ([PR#884](https://github.com/Scalingo/cli/pull/884))
* feat(deployments): add Image Size to the list of deployments ([PR#894](https://github.com/Scalingo/cli/pull/894))
* fix(backups): backup flag is not required ([PR#892](https://github.com/Scalingo/cli/pull/892))
* build(publish): replace `rm-dist` with `clean` ([PR#893](https://github.com/Scalingo/cli/pull/893))
* feat(review-apps): add an option to manage review app creation from forks ([PR#882](https://github.com/Scalingo/cli/pull/882))
* chore(deps): update dependencies
  * golang.org/x/net from 0.5.0 to 0.7.0
  * github.com/stretchr/testify from 1.8.1 to 1.8.2
  * golang.org/x/mod from 0.7.0 to 0.8.0
  * github.com/pjbgf/sha1cd from 0.2.3 to 0.3.0
  * github.com/golang-jwt/jwt/v4 from 4.4.3 to 4.5.0
  * github.com/go-git/go-billy/v5 from 5.4.0 to 5.4.1

### 1.27.2

* bump: go-scalingo v6.2.0
* fix(commands): cli help refactoring
* fix(commands): brought back --help command
* fix: use a current stacks in help

### 1.27.1

* fix(cli-build) Compile the CLI statically to prevent GLIBC incompatibility [#863](https://github.com/Scalingo/cli/pull/863)

### 1.27.0

* fix(commands): display help usage message if needed [#835](https://github.com/Scalingo/cli/pull/835)
* feat(logs): allow DB type aliases [#829](https://github.com/Scalingo/cli/pull/829)
* feat(db-console): add `influx-console` and `mongodb-console` aliases [#830](https://github.com/Scalingo/cli/pull/830)
* feat(send-signal): `send-signal` command added [#833](https://github.com/Scalingo/cli/pull/833)
* feat(send-signal): Allow to send signals to multiple containers via its type [#841](https://github.com/Scalingo/cli/pull/841)

### 1.26.1

* deps(urfave/cli): Bumpd from 2.17.1 to 2.23.5: fix options display regression in --help [#814](https://github.com/Scalingo/cli/pull/814)
* deps(websocket): Replace golang.org/x/net/websocket with github.com/gorilla/websocket [#822](https://github.com/Scalingo/cli/pull/822)
* deps(survey): upgrade gopkg.in/AlecAivazis/survey.v1 to github.com/AlecAivazis/survey/v2 [#821](https://github.com/Scalingo/cli/pull/821)
* fix(logdrains): delete Logentries mentions [#820](https://github.com/Scalingo/cli/pull/820)

### 1.26.0

* feat(db-commands) Add `database-enable-feature` and `database-disable-feature` for database addons [#807](https://github.com/Scalingo/cli/pull/807)
* deps(go-scalingo) Bump from 5.2.2 to 6.0.0 [#804](https://github.com/Scalingo/cli/pull/804)
* regression(cmd-alias) 'scale' command can again be invoked by 's' [#802](https://github.com/Scalingo/cli/pull/802)
* regression(cmd-alias) 'run' command can again be invoked by 'r' [#802](https://github.com/Scalingo/cli/pull/802)

### 1.25.1

This release is mandatory to keep the `logs` command working in the coming weeks.

* feat(logs): add support for the new authentication mechanism

### 1.25.0

#### Changed

* [BREAKING] chore(deps): Upgrade urfave/cli to v2.16 Command arguments must now come after the command flags.
  For example, `scalingo --app my-app integration-link-create https://ghe.example.org/test/test-app --auto-deploy`
  must be rewritten `scalingo --app my-app integration-link-create --auto-deploy https://ghe.example.org/test/test-app`
  [#774](https://github.com/Scalingo/cli/pull/774)

* chore(deps): bump github.com/Scalingo/go-scalingo from v4.16 to v5.2 [#775](https://github.com/Scalingo/cli/pull/775)

#### Added

* feat(stacks): hide deprecated stacks when listing them with `scalingo stacks`.
  It is still possible to list deprecated stacks and to show deprecation dates
  with `scalingo stacks --with-deprecated` [#776](https://github.com/Scalingo/cli/pull/776)

### 1.24.2

* fix(releases): Correctly auto generate changelog [#772](https://github.com/Scalingo/cli/pull/772)
* feat(changelog): add changelog command and show the changelog after an update [#767](https://github.com/Scalingo/cli/pull/767)

### 1.24.1

* fix(releases): Windows archives should be in the zip format [#769](https://github.com/Scalingo/cli/pull/769)

### 1.24.0

* feat(db): environment variable name is customizable [#759](https://github.com/Scalingo/cli/issues/759)
* feat(goreleaser): use goreleaser to make releases using github action [#752](https://github.com/Scalingo/cli/issues/752)
* feat(make-release): use gox [#747](https://github.com/Scalingo/cli/pull/747)
* fix(install.sh): better error message if fails to get the version [#748](https://github.com/Scalingo/cli/pull/748)
* feat(logs): possibility to get the addon logs by using its type (e.g. Redis) [#745](https://github.com/Scalingo/cli/pull/745)
* fix(error): interpret raw newline [#744](https://github.com/Scalingo/cli/pull/744)
* fix(stats): display memory as IEC size [#742](https://github.com/Scalingo/cli/pull/742)
* fix(autoscalers): min containers cannot be 1 [#741](https://github.com/Scalingo/cli/pull/741)
* fix(appdetect): Do not crash if there is a remote without any URL [#755](https://github.com/Scalingo/cli/pull/755)
* refactor(run_unix): replace github.com/heroku/hk/term with golang.org/x/term [#740](https://github.com/Scalingo/cli/pull/740)
* build(deps): bump github.com/stretchr/testify from 1.7.2 to 1.8.0
* build(deps): bump github.com/cheggaaa/pb/v3 from 3.0.8 to 3.1.0
* build(deps): bump github.com/briandowns/spinner from 1.18.1 to 1.19.0

### 1.23.0

* fix(domains): replace DomainsUpdate with Domain\*Certificate [#731](https://github.com/Scalingo/cli/pull/731)
* chore(go): use go 1.17 [#728](https://github.com/Scalingo/cli/pull/728)
* feat(app): Region detection from Git remote [#724](https://github.com/Scalingo/cli/pull/724)
* feat(app): add a flag --force to destroy an app without interactive confirmation [#721](https://github.com/Scalingo/cli/pull/721)
* build(deps): bump github.com/briandowns/spinner from 1.18.0 to 1.18.1
* build(deps): bump github.com/stretchr/testify from 1.7.0 to 1.7.1
* build(deps): bump github.com/urfave/cli from 1.22.5 to 1.22.8
* build(deps): bump github.com/Scalingo/go-utils/errors from 1.1.0 to 1.1.1
* build(deps): bump github.com/Scalingo/go-utils/retry from 1.1.0 to 1.1.1

### 1.22.1

* feat(logs): allow to filter logs to only show router logs [#707](https://github.com/Scalingo/cli/pull/707)
* fix(install): disable CLI update checker [#710](https://github.com/Scalingo/cli/pull/710)
* fix(run): display an error message for detached one-off with uploaded files [#712](https://github.com/Scalingo/cli/pull/712)
* fix(run): --region flag in help message [#713](https://github.com/Scalingo/cli/pull/713)
* build(deps): bump github.com/cheggaaa/pb from 1.0.29 to 3.0.8
* build(deps): bump github.com/briandowns/spinner from 1.16.0 to 1.18.0

### 1.22.0

* feat(logs-archives): add logs archives for addons [#694](https://github.com/Scalingo/cli/pull/694)
* feat(pgsql-console): add `psql-console` and `postgresql-console` aliases for `pgsql-console` command and replace duplicated commands with aliases [#693](https://github.com/Scalingo/cli/pull/693)
* feat(router-logs): add command `router-logs` to enable/disable router logs on your application [#692](https://github.com/Scalingo/cli/pull/692)
* feat(open): add command `open` to open app on default browser [#691](https://github.com/Scalingo/cli/pull/691)
* feat(addons-info): add command `addons-info` to display information of an add-on [#689](https://github.com/Scalingo/cli/pull/689)
* feat(dashboard): add command `dashboard` to open dashboard of specified app on default browser [#686](https://github.com/Scalingo/cli/pull/686)
* chore(deps): replace github.com/howeyc/gopass with golang.org/x/term [#703](https://github.com/Scalingo/cli/pull/703)
* fix(update): change data stream on which warning is displayed from stdout to stderr [#698](https://github.com/Scalingo/cli/pull/698)
* build(deps): bump github.com/Scalingo/go-utils/errors from 1.0.0 to 1.1.0 [#687](https://github.com/Scalingo/cli/pull/687)
* build(deps): bump github.com/Scalingo/go-utils/retry from 1.0.0 to 1.1.0 [#688](https://github.com/Scalingo/cli/pull/688)

### 1.21.2

* feat(cron-task): add fields to cron tasks list [#685](https://github.com/Scalingo/cli/pull/685)

### 1.21.1

* feat(deploy): add --no-follow flag to detach deployment logs [#679](https://github.com/Scalingo/cli/pull/679)
* build(deps): bump github.com/fatih/color from 1.12.0 to 1.13.0 [#680](https://github.com/Scalingo/cli/pull/680)
* fix(events): fix display for events new_autoscaler, edit_autoscaler and repeated_crash in timeline [#684](https://github.com/Scalingo/cli/pull/684)

### 1.21.0

* fix(backup-download): Now download the last successful backup [#663](https://github.com/Scalingo/cli/pull/663)
* fix(review-app-url): add app url in the table [#664](https://github.com/Scalingo/cli/pull/664)
* fix(backups-download): print error on stderr if the backup status is not 'done' [#665](https://github.com/Scalingo/cli/pull/665)
* feat(command): add command `cron-tasks` to list cron tasks of an application [#670](https://github.com/Scalingo/cli/pull/670)
* fix(run): allow an equals sign (=) in environment variable values [#674](https://github.com/Scalingo/cli/pull/674)
* bump github.com/fatih/color from 1.11.0 to 1.12.0: add support for `NO_COLOR` environment variable.
* update go version to 1.16 and replace ioutil by io/os [#666](https://github.com/Scalingo/cli/pull/666)
* bump github.com/briandowns/spinner from 1.12.0 to 1.16.0
* bump github.com/go-git/go-git/v5 from 5.4.1 to 5.4.2
* bump github.com/Scalingo/go-scalingo/v4 from 4.13.1 to 4.14.0

### 1.20.2

* fix(logs) Buffer logs to prevent requests from timing out [#659](https://github.com/Scalingo/cli/pull/659)
* feat(integration-link-create): better catch 401 from SCM API [#651](https://github.com/Scalingo/cli/pull/651)
* bump github.com/fatih/color from 1.10.0 to 1.11.0 [#654](https://github.com/Scalingo/cli/pull/654)

### 1.20.1

* fix(backup_download): no log on stdout when using '--output -' flag [#646](https://github.com/Scalingo/cli/pull/646)
* `deployment-logs` displays the last deployment logs if none is provided [#647](https://github.com/Scalingo/cli/pull/647)

### 1.20.0

* Add `env-get` command to retrieve the value of a specific environment variable [#643](https://github.com/Scalingo/cli/pull/643)
* Error message are outputted on stderr [#639](https://github.com/Scalingo/cli/pull/639)
* Automatically prefix the integration URL with https if none is provided [#642](https://github.com/Scalingo/cli/pull/642)
* `backups-download` downloads the most recent backup if none is specified [#636](https://github.com/Scalingo/cli/pull/636)
* Add `deployment-cache-delete` as an alias to `deployment-delete-cache` [#635](https://github.com/Scalingo/cli/pull/635)
* `one-off-stop` command to stop a running one-off container [#633](https://github.com/Scalingo/cli/pull/633)
* `ps` command returns the list of application's containers [#632](https://github.com/Scalingo/cli/pull/632)
* `ps` command is renamed `scale` [#631](https://github.com/Scalingo/cli/pull/631)
* Show the addon status on `scalingo addons` [#604](https://github.com/Scalingo/cli/pull/604)
* Add informative error in case of container type error when scaling an application
  [602](https://github.com/Scalingo/cli/pull/602)
* Update various dependencies:
  * github.com/fatih/color
  * github.com/briandowns/spinner
  * github.com/gosuri/uilive
  * github.com/stretchr/testify
  * github.com/urfave/cli
  * gopkg.in/AlecAivazis/survey.v1
  * github.com/cheggaaa/pb
  * gopkg.in/src-d/go-git.v4
  * github.com/Scalingo/go-scalingo/v4
  * github.com/ScaleFT/sshkeys
  * github.com/heroku/hk
  * github.com/howeyc/gopass
  * github.com/olekukonko/tablewriter
  * golang.org/x/crypto
  * golang.org/x/net
* Update Scalingo internal dependencies to the Go Modules version [#613](https://github.com/Scalingo/cli/pull/613)

### 1.19.3

* Add support for new log drain/addon log drain events
  [#600](https://github.com/Scalingo/cli/pull/600)

### 1.19.2

* Use an unauthenticated to fetch region cache on login
  [#599](https://github.com/Scalingo/cli/pull/599)

### 1.19.1

* Ask user to end its free trial if asking for automatic deployment of review
  apps [#588](https://github.com/Scalingo/cli/pull/588)
* Fetch the SSH Login endpoint from our region metadata
  [#592](https://github.com/Scalingo/cli/pull/592)
* Fix a panic when log-drains-remove was launched without any args
  [#595](https://github.com/Scalingo/cli/pull/595)
* Add support for ECDSA SSH keys

### 1.19.0

* Handle new `queued` deployment status
  [#586](https://github.com/Scalingo/cli/pull/586)
  [go-scalingo #177](https://github.com/Scalingo/go-scalingo/pull/177)
* Remove a log drain from an addon
  [#577](https://github.com/Scalingo/cli/pull/577)
  [go-scalingo #175](https://github.com/Scalingo/go-scalingo/pull/175)
* Add a log drain to an addon
  [#575](https://github.com/Scalingo/cli/pull/575)
  [go-scalingo #174](https://github.com/Scalingo/go-scalingo/pull/174)
* Add list log drains by addons command
  [#572](https://github.com/Scalingo/cli/pull/572)
  [go-scalingo #172](https://github.com/Scalingo/go-scalingo/pull/172)
* Update the `apps` command to display the application status
  [#582](https://github.com/Scalingo/cli/pull/582)

### 1.18.1

* Fix infinite recursion on some events with `scalingo timeline` [#570](#570)
* Fix log drain response parsing [#570](#570)

### 1.18.0

* Add log drains list command
  [#560](https://github.com/Scalingo/cli/pull/560)
  [go-scalingo #163](https://github.com/Scalingo/go-scalingo/pull/163)
* Add log drains add command
  [#561](https://github.com/Scalingo/cli/pull/561)
  [go-scalingo #164](https://github.com/Scalingo/go-scalingo/pull/164)
* Add log drains delete command
  [#567](https://github.com/Scalingo/cli/pull/567)
  [go-scalingo #167](https://github.com/Scalingo/go-scalingo/pull/167)

### 1.17.0

* Region Migration: Better output for the `migration-follow` command to explain what is happening [#549](https://github.com/Scalingo/cli/pull/549)
* Region Migration: Retry if it fails to refresh the migration status [#550](https://github.com/Scalingo/cli/pull/550)
* Region Migration: Add instructions to change local Git URL at the end of the migration [#551](https://github.com/Scalingo/cli/pull/551)
* Deployment: Send a non 0 exit code when a deployment fails [#563](https://github.com/Scalingo/cli/pull/563)
* Bugfix: fix support for `addon_updated` and `start_region_migration` event type [#558](https://github.com/Scalingo/cli/pull/558)
* Bugfix: fix author of `restart` and `edit_variable` when it's an addon [#558](https://github.com/Scalingo/cli/pull/558)

### 1.16.8

* Better output for migration commands
  [#540](https://github.com/Scalingo/cli/pull/540)
  [#542](https://github.com/Scalingo/cli/pull/542)

### 1.16.7

* Improve error message if unknown app, suggests to try on a different region
  [#524](https://github.com/Scalingo/cli/pull/524)
* Add option `--force` to the command `git-setup`
  [#527](https://github.com/Scalingo/cli/pull/527)
* Correctly display error messages from the API
  [#528](https://github.com/Scalingo/cli/pull/528)

### 1.16.6

* Bugfix: integration-link-create: command will always fail if there's already a link [#520](https://github.com/Scalingo/cli/issues/520)

### 1.16.5

* Add `--bind` arg to `db-tunnel` command that let us bind a custom host (and not only 127.0.0.1) [#517](https://github.com/Scalingo/cli/pull/517)
* Error message if unknown app suggests to try on a different region [#519](https://github.com/Scalingo/cli/pull/519)

* Bugfix: Encrypted key with new OpenSSH format header for private keys was broken [#513](https://github.com/Scalingo/cli/pull/513)
* Bugfix: db-tunnel: better handling of short lived connections  [#517](https://github.com/Scalingo/cli/pull/517)

### 1.16.4

* Bugfix: Unable to create integration link for an app [#510](https://github.com/Scalingo/cli/pull/510)

### 1.16.3

* Display a better error message when app and addon creation are disabled on a
  region [#503](https://github.com/Scalingo/cli/pull/503),
  [#504](https://github.com/Scalingo/cli/pull/504)

### 1.16.2

* Do not get the SCM integration when showing the integration link
  [#492](https://github.com/Scalingo/cli/pull/492)
* Better handling of 404 from auth query
  [#493](https://github.com/Scalingo/cli/pull/493)

### 1.16.1

* Global commands do not need to query the list of regions
  [#487](https://github.com/Scalingo/cli/pull/487)

### 1.16.0

* Add `integrations`, `integrations-add`, `integrations-delete` and
  `integrations-import-keys` commands
  [#444](https://github.com/Scalingo/cli/pull/444)
* Add `integration-link`, `integration-link-create`, `integration-link-update`,
  `integration-link-delete`, `integration-link-manual-deploy`,
  `integration-link-manual-review-app` commands
  [#458](https://github.com/Scalingo/cli/pull/458)
* Add support for new SCM-related events
  [#467](https://github.com/Scalingo/cli/pull/467) and
  [#458](https://github.com/Scalingo/cli/pull/458)
* Bugfix: Do not disconnect user if the API returns 401
  [#463](https://github.com/Scalingo/cli/pull/463)
* Add duration to the deployments list
  [#477](https://github.com/Scalingo/cli/pull/477)
* Add support for new_user event [#473](https://github.com/Scalingo/cli/pull/473)
* [Redis console] Better error message if TLS connections are enforced
  [#480](https://github.com/Scalingo/cli/pull/480)
* Default "Destroy on Stale" for interactive integration link creation is "No"
  [#485](https://github.com/Scalingo/cli/pull/485)

### 1.15.1

* Broken retro-compatibility of the `backup-download` command [#460](https://github.com/Scalingo/cli/pull/460)

### 1.15.0

* Bugfix if notifier send all events [#453](https://github.com/Scalingo/cli/pull/453)
* Only fill regions cache for regional commands [#449](https://github.com/Scalingo/cli/pull/449)
* Add `backups-create` command [#452](https://github.com/Scalingo/cli/pull/452)
* Deprecate `backup-download` in favor of `backups-download` [#452](https://github.com/Scalingo/cli/pull/452)
* Add `migrations`, `migrations-create` and `migrations-follow` [#446](https://github.com/Scalingo/cli/pull/446)
* Add periodic backups configuration commands [#455](https://github.com/Scalingo/cli/pull/455)
* PeriodicBackupsScheduledAt can contain multiple values [#457](https://github.com/Scalingo/cli/pull/457)

### 1.14.1

* Fix Git remote URL parsing to detect app name on Outscale [#443](https://github.com/Scalingo/cli/pull/443)

### 1.14.0

* Use user's default region [#411](https://github.com/Scalingo/cli/pull/441)

### 1.13.0

* Add the `apps-info` command [#438](https://github.com/Scalingo/cli/pull/438)
* Display request ID in debug logs [#435](https://github.com/Scalingo/cli/pull/435)
* Add `git-setup` and `git-show` commands [#431](https://github.com/Scalingo/cli/pull/431)
* Remove dependency to an old Git lib for a more battle tested one [#434](https://github.com/Scalingo/cli/pull/434)

### 1.12.0

* Initial support for multi-region [#425](https://github.com/Scalingo/cli/pull/425)
* [self] Check logged in user with the API [#427](https://github.com/Scalingo/cli/pull/427)

### 1.11.0

* [alerts] Add support for the duration_before_trigger attribute [#407](https://github.com/Scalingo/cli/pull/407)
* [domains] Handle Let's Encrypt certificate status [#410](https://github.com/Scalingo/cli/pull/410)
* Fix login with SSH [#419](https://github.com/Scalingo/cli/pull/419)

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
* <leo@scalingo.com> <Leo Unbekandt>
* <mail2tevin@gmail.com> <Tevin Zhang>

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
