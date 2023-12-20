# Changelog

## To Be Released

## 6.7.5

fix(events): `link_scm` data types

## 6.7.4

* feat(database): add user management

## 6.7.3

* feat(events): add maintenance events
* refactor(tokens): deprecate `ErrOTPRequired` and `IsOTPRequired`
* fix(events): disambiguate event name when restarting containers
* fix(events): app run event wording

## 6.7.2

* fix(maintenance): change listing pagination meta to standard pagination meta

## 6.7.1

* style(readme): update api version

## 6.7.0

* feat(maintenance): add maintenance windows
* feat(maintenance): add maintenance listing
* feat(maintenance): add maintenance show

## 6.6.0

* feat(scm-repo-link): fetch a Pull Request data

## 6.5.0

* feat(scm-repo-link): add ManualDeployResponse

## 6.4.0

* feat(one-off): allow attached one-off to be run asynchronously
* feat(stacks): add support for Deployment `stack_base_image`
* feat(backup): add `StartedAt` and `Method` fields

## 6.3.0

* feat(apps): add support for `hds_resource`
* feat(deployments): add support for Deployment `image_size`
* feat(events): Stack changed event

BREAKING CHANGE:

* Remove the GithubLinkService in favor of SCMRepoLinkService

## 6.2.0

* feat(status code): handle too many requests exception

## 6.1.0

* feat(events): Add Detached field for EventRun [#292](https://github.com/Scalingo/go-scalingo/pull/292)
* chore(deps): bump github.com/golang-jwt/jwt/v4 from 4.4.2 to 4.4.3
* feat(send-signal): new function `ContainersKill` to request the send-signal api's endpoint [#295](https://github.com/Scalingo/go-scalingo/pull/295)

## 6.0.1

* chore(deps): bump github.com/stretchr/testify from 1.8.0 to 1.8.1
* chore(deps): replace golang.org/x/net/websocket -> github.com/gorilla/websocket

## 6.0.0

BREAKING CHANGES:

Linting:
* `App.GitUrl` → `App.GitURL`
* `ContainerStat.CpuUsage` → `ContainerStat.CPUUsage`
* `LogsArchiveItem.Url` → `LogsArchiveItem.URL`
* `PullRequest.Url` → `PullRequest.URL`
* `PullRequest.HtmlUrl` → `PullRequest.HTMLURL`

Making things more homogeneous:
* `PeriodicBackupsConfig` → `DatabaseUpdatePeriodicBackupsConfig`
* `PeriodicBackupsConfigParams` → `DatabaseUpdatePeriodicBackupsConfigParams`

Changelog:

* feat(data_access_consent): Add DataAccessConsent field on `app` [#287](https://github.com/Scalingo/go-scalingo/pull/287)
* refactor(linting): Fix linter offense on App, ContainerStat, LogsArchiveItem and PullRequest public structs [#288](https://github.com/Scalingo/go-scalingo/pull/288)
* feat(databases): Add mock of DatabasesService, unifomize naming [#286](https://github.com/Scalingo/go-scalingo/pull/286)
* feat(databases): Add ability to add/remove features [#285](https://github.com/Scalingo/go-scalingo/pull/285)

## 5.2.2

*  feat(logs): add timestamp to the logs URL [#283](https://github.com/Scalingo/go-scalingo/pull/283)

## 5.2.1

* fix(notifiers): add missing fields in `Notifier` [#278](https://github.com/Scalingo/go-scalingo/pull/278)
* fix(scm_repo_link): add missing field `URL` in the structure `SCMRepoLink` [#281](https://github.com/Scalingo/go-scalingo/pull/282)

## 5.2.0

* feat(alerts): add missing fields for alerts endpoint [#276](https://github.com/Scalingo/go-scalingo/pull/276)
* feat(invoices): add invoices endpoint [#269](https://github.com/Scalingo/go-scalingo/pull/269)
* feat(scm integrations): add missing fields for scm integrations endpoint [#270](https://github.com/Scalingo/go-scalingo/pull/271).
  The structure `SCMIntegration` embraces the API return.
  `Owner` structure has been shared.
* fix(addon provider): add missing fields for addon providers endpoint [#272](https://github.com/Scalingo/go-scalingo/pull/272).
  The structure `Plan` and `AddonProvider` embrace the API return.
  `Category` structure has been added.

## 5.1.0

* feat(stacks): Add deprecation date to stacks in `DeprecatedAt` attribute, and
  a `IsDeprecated` method to tell if a stack is deprecated or not based on its
  deprecation date [#265](https://github.com/Scalingo/go-scalingo/pull/265)
* feat(new-events): Add compatibility with new events related to Account
  management and Database feature management
  [#252](https://github.com/Scalingo/go-scalingo/pull/252)
* dev(generators): Use Go generators to generate Events boilerplate
  [#252](https://github.com/Scalingo/go-scalingo/pull/252)

## 5.0.0

BREAKING CHANGES:

* Added a `context.Context` argument to the methods
* Deprecated the following methods:
  * `AppsPs`: replaced with `AppsContainerTypes`
  * `Container.Size`: replaced with `ContainerSize`
  * `DomainsUpdate`: replaced with `DomainsSetCanonical`, `DomainUnsetCanonical`, `DomainSetCertificate` or `DomainUnsetCertificate`
  * `Login`: replaced with the OAuth flow
  * `StopFreeTrial`: replaced with `UserStopFreeTrial` method
  * Make `FillDefaultValues` private: it shouldn't have been used.

Changelog:

* feat: remove deprecated methods [#263](https://github.com/Scalingo/go-scalingo/pull/263)
* feat(http): add X-Request-ID if ID is in the context [#262](https://github.com/Scalingo/go-scalingo/pull/262)
* refactor(events): moved alert related events to events_addon.go [#264](https://github.com/Scalingo/go-scalingo/pull/264)
* refactor(events): moved addons related events to events_addon.go [#261](https://github.com/Scalingo/go-scalingo/pull/261)
* refactor(events): moved direct app related events to events_app.go [#260](https://github.com/Scalingo/go-scalingo/pull/260)
* feat(events): implement two factor auth related events [#259](https://github.com/Scalingo/go-scalingo/pull/259)
* feat(events): implement token related events [#257](https://github.com/Scalingo/go-scalingo/pull/257)
* feat(events): implement event for the update of a key [#256](https://github.com/Scalingo/go-scalingo/pull/256)
* feat(events): implement event for the update of an hds contact [#255](https://github.com/Scalingo/go-scalingo/pull/255)
* feat(events): implement event for the creation of a data access consent [#255](https://github.com/Scalingo/go-scalingo/pull/255)
* build(deps): bump github.com/stretchr/testify from 1.7.2 to 1.8.0
* build(deps): bump github.com/golang-jwt/jwt/v4 from 4.4.1 to 4.4.2

## 4.16.0

* refactor: replace use of the deprecated ioutil package [#238](https://github.com/Scalingo/go-scalingo/pull/238)
* feat(scm-repo-link): Add `SCMRepoLinkList` method [#241](https://github.com/Scalingo/go-scalingo/pull/241) and [#245](https://github.com/Scalingo/go-scalingo/pull/245)
* build(deps): bump github.com/golang-jwt/jwt/v4 from 4.1.0 to 4.4.1
* build(deps): bump github.com/stretchr/testify from 1.7.0 to 1.7.2
* feat(log-drains): The `LogDrainsAddonList` now returns the list of log drains [#246](https://github.com/Scalingo/go-scalingo/pull/246)
* feat(log-drains): Cleanup the LogDrain struct [#246](https://github.com/Scalingo/go-scalingo/pull/246)
* refactor(domain): deprecate DomainsUpdate [#250](https://github.com/Scalingo/go-scalingo/pull/250)
* feat(domains): add DomainSetCertificate and DomainUnsetCertificate methods [#250](https://github.com/Scalingo/go-scalingo/pull/250)

## 4.15.1

* chore: generate missing mocks from v4.15.0

## 4.15.0

* feat: add support to list container sizes [#234](https://github.com/Scalingo/go-scalingo/pull/234)
* feat(addons): add the `AddonLogsArchives` method [#235](https://github.com/Scalingo/go-scalingo/pull/235)
* feat(apps): add `AppsRouterLogs` method [#236](https://github.com/Scalingo/go-scalingo/pull/236)

## 4.14.3

* feat(cron-task): Add fields to cron tasks model [#231](https://github.com/Scalingo/go-scalingo/pull/231)

## 4.14.2

fix(events): change type expected in json parsing from implicit int to explicit string [#230](https://github.com/Scalingo/go-scalingo/pull/230)

## 4.14.1

* build(deps): bump github.com/golang-jwt/jwt/v4 from 4.0.0 to 4.1.0 [#225](https://github.com/Scalingo/go-scalingo/pull/225)
* chore: use jwt.RegistredClaims instead of jwt.StandardClaims [#227](https://github.com/Scalingo/go-scalingo/pull/227)
* feat(events): implement event for the edition of an autoscaler. [#228](https://github.com/Scalingo/go-scalingo/pull/228)
* feat(events): implement event for repeated crash [#229](https://github.com/Scalingo/go-scalingo/pull/229)

## 4.14.0

* fix(apps): typo on the `UpdatedAt` field of the Apps [#220](https://github.com/Scalingo/go-scalingo/pull/220)
* fix(token): correctly deserialize the `ID` [#222](https://github.com/Scalingo/go-scalingo/pull/222)
* feat(token): add the `LastUsedAt` field [#222](https://github.com/Scalingo/go-scalingo/pull/222)
* feat(deployments): add the `DeploymentListWithPagination` method [#221](https://github.com/Scalingo/go-scalingo/pull/221)
* feat(addon): add the `ProvisionedAt` and `DeprovisionedAt` fields [#224](https://github.com/Scalingo/go-scalingo/pull/224)

## 4.13.1

* build(deps): bump github.com/golang/mock from 1.5.0 to 1.6.0
* Add `DeploymentUUID` to the deployment event and fix a typo [#216](https://github.com/Scalingo/go-scalingo/pull/216)
* build(deps): replace deprecated github.com/dgrijalva/jwt-go with github.com/golang-jwt/jwt
* chore: update from Go 1.13 to 1.17

## 4.13.0

* feat(domain): add SSL status [#207](https://github.com/Scalingo/go-scalingo/pull/207)
* feat: add user agent to auth and DB API requests [#206](https://github.com/Scalingo/go-scalingo/pull/206)
* feat(event): add support for notifier events [#205](https://github.com/Scalingo/go-scalingo/pull/205)
* Add `ContainersStop` to stop a running one-off container [#204](https://github.com/Scalingo/go-scalingo/pull/204)
* Add `AppsContainersPs` to list all application's containers [#204](https://github.com/Scalingo/go-scalingo/pull/204) and [#208](https://github.com/Scalingo/go-scalingo/pull/208)
* Add `AppsContainerTypes` and deprecate `AppsPs` [#204](https://github.com/Scalingo/go-scalingo/pull/204) and [#208](https://github.com/Scalingo/go-scalingo/pull/208)
* Add support for route `GET /apps/{app_id}/cron_tasks` [#214](https://github.com/Scalingo/go-scalingo/pull/214)

## 4.12.0

* Add `Limits` field to `App` [#203](https://github.com/Scalingo/go-scalingo/pull/203)

## 4.11.1

* Rename `DeployEvent*` structures to `DeploymentEvent*` [#201](https://github.com/Scalingo/go-scalingo/pull/201)

## 4.11.0

* Add `DeployEvent` structure which represents a deployment stream event sent on the websocket [#200](https://github.com/Scalingo/go-scalingo/pull/200)

## 4.10.0

* Remove headers from `Request` of `http.parseJSON` returned error
  [#199](https://github.com/Scalingo/go-scalingo/pull/199)
* Add missing fields `Plan.SKU` (corresponding catalogue name)
  [#198](https://github.com/Scalingo/go-scalingo/pull/198) thanks @mackwic

## 4.9.0

* Add `Flags` field to `App`

## 4.8.2

* Update mocks

## 4.8.1

* Add support for database_type_versions#show endpoint
* Go module v4

## 4.8.0

* Make the API prefix version configurable for all APIs
* Deserialize app.Region returned by API
* Fix panic GithubLinkShow when an app does not have one
* Migrate to go mod

## 4.7.2

* Add EventNewLogDrainType and EventDeleteLogDrainType event
* Add EventNewAddonLogDrainType and EventDeleteAddonLogDrainType event

## 4.7.1

* Add support for optionally authenticated routes

## 4.7.0

* Add support for route `PATCH /users/stop_free_trial`
* Deprecate the use of `UpdateUser` to stop a user free trial

## 4.6.0

* Add `queued` deployment status and add it to `IsFinishedString`

## v4.5.8

* Add `LogDrainAddonRemove` function to remove a log drain from an addon

## v4.5.7

* Add `LogDrainAddonAdd` function to add a log drain to an addon

## v4.5.6

* Add `LogDrainsAddonList` function to list log drains of an addon

## v4.5.5

* Update GET /log_drains route to return object instead of array

## v4.5.4

* Fix infinite recursion on some events

## v4.5.3

* Add `LogDrainRemove` function to remove log drain from an application

## v4.5.2

* Add `LogDrainsAdd` function to add log drains to an application

## v4.5.1

* Add `HasFailed` and `HasFailedString` to deployments

## v4.5.0

* Add `LogDrainsService` service
* Add `LogDrainsList` function to list log drains of an application

## v4.4.1

* Fix support for `addon_updated` and `start_region_migration` event types
* Fix who created a `restart` or `edit_variable` event types

## v4.4.0

* Rename `scheduled` to `created` in RegionMigrationStatus list
* Add `aborting` status to the list of RegionMigrationStatus

## v4.3.1

* Add `Source` field to RegionMigration

## v4.3.0

* Add StartRegionMigration event

## v4.2.0

* Add a JWT token cache when a client is reused, a new JWT token will be only used if the previous has expired.

## v4.1.0

* Change AddonUpgrade signature, take a struct instead of a simple string for planID

## v4.0.0

* Add the ability to pass options when provisioning an addon

## v3.2.0

* Add ability to reach AuthService with the `StaticTokenGenerator` (Use of a client with a predefined JWT)

## v3.1.1

* Final structure of `AddonUpdatedEvent`
  * Add `AddonResourceID` `AddonPlanName` and `AddonProviderName`

## v3.1.0

* New Event `AddonUpdatedEvent`

* [API CHANGE] Remove `NewNotification`/`EditNotification`/`DeleteNotification` events
* [API CHANGE] Remove interface NotificationsService and all included methods.
```
  NotificationsList(app string) ([]*Notification, error)
  NotificationProvision(app, webHookURL string) (NotificationRes, error)
  NotificationUpdate(app, ID, webHookURL string) (NotificationRes, error)
  NotificationDestroy(app, ID string) error
```
  From now `NotifiersService` and methods should be used instead.
* [API CHANGE] Remove `SCMIntegrationUUID` field from `SCMRepoLinkCreateParams` and `SCMRepoLink` structs.
  [#148](https://github.com/Scalingo/go-scalingo/pull/148)

## v3.0.8

* Add ability to give DstAppName to a RegionMigration
  [#145](https://github.com/Scalingo/go-scalingo/pull/145)

## v3.0.7

* Improve error management, get Code field from API for 403 and 400 errors
  [#144](https://github.com/Scalingo/go-scalingo/pull/144)

## v3.0.6

* Add error management for 403 forbidden errors
  [#143](https://github.com/Scalingo/go-scalingo/pull/143)

## v3.0.5

* Add `SCMType` to `SCMRepoLink`
  [#138](https://github.com/Scalingo/go-scalingo/pull/138)
* Better handling of 404 from API
  [#140](https://github.com/Scalingo/go-scalingo/pull/140)
* Display the user agent in debug output
  [#139](https://github.com/Scalingo/go-scalingo/pull/139)

## v3.0.4

* Removed the `SCMRepoLinkParams` struct. You should use
  `SCMRepoLinkCreateParams` instead for creation call
* Added the `SCMRepoLinkUpdateParams` struct, you should use this for update call.

## v3.0.3

* Add support (again) for the new SCM events [#135](https://github.com/Scalingo/go-scalingo/pull/135)
* Add support for the new user event [#136](https://github.com/Scalingo/go-scalingo/pull/136)

## v3.0.2

* Add support for the new SCM events [#134](https://github.com/Scalingo/go-scalingo/pull/134)

## v3.0.1

* Fix HTTP verb for Repolink update
* Set `omitempty` on ParentID and StackID when creating an App

## v3.0.0

* Add region support [#129](https://github.com/Scalingo/go-scalingo/pull/129)
This comes with the following changes:
 - Removed the `NewClient` method. You should use `New` instead
 - Add the `Region` field to `ClientConfig`
 - Drop support for the `SCALINGO_AUTH_URL`, `SCALINGO_API_URL`, `SCALINGO_DB_URL` environment variables. You should use the corresponding `ClientConfig` fields instead.
 - Remove the default values for `APIEndpoint` and `DatabaseAPIEndpoint`
 - Add some basic examples

## v2.5.5

* Add SCMIntegrationService [#126](https://github.com/Scalingo/go-scalingo/pull/126)

## v2.5.4

* Add constants for SCM integrations name [#124](https://github.com/Scalingo/go-scalingo/pull/124)

## v2.5.3

* PeriodicBackupsScheduledAt is a slice of int

## v2.5.2

* `sand_ip` is renamed to `private_ip`
* Database#Show is now correctly parsing JSON from database API

## v2.5.1

* Add `DatabaseShow`

## v2.5.0

* Add periodic backups configuration method
* Add repo link calls
* Backup creation returns a backup

## v2.4.9

* Add `RegionMigrationsService`
* Add `BackupCreate` method

## v2.4.8

* Avoid using `errgo.Mask`

## v2.4.7

* Add StackID field when creating a new app

## v2.4.6

* Add Default attribute to Region

## v2.4.5

* Add missing fields in the App struct

## v2.4.4

* Display request ID in the debug logs

## v2.4.3

* Add SSH to Region to configure SSH endpoint for SSH-based operations

## v2.4.2

* UsersSelf is now based on AuthenticationService not API anymore
* Add RegionsList() to get the list of available platform regions

## v2.4.1

* Implement new methods on Client: `EventTypesList` and `EventCategoriesList`

## v2.4.0

* Update `Notifier` methods to match the API: ie. accept SelectedEventIDs as input and output

## v2.3.0

* Update `NotifierUpdate`/`NotifierProvision` params, Add missing field `Notifier.SendAllAlerts`

## v2.2.1

* Add `Description` and `LogoURL` to `NotificationPlatform` struct

## v2.2.0

* Add `OperationsShowFromURL(url string) (*Operation, error)` to ease the use
  of the Operation URL returned after a Scale/Restart action
* Add `OperationStatus` and `OperationType` types with the right constants
* Remove `Plan.TextDescription` which was never used

## v2.1.1

* Add StaticTokenGenerator in ClientConfig to ensure retrocompatibility

## v2.1.0

* StacksList() to list available runtime stacks
* Add AppsSetStack() to update the stack of an app

## v2.0.0

* Integration of Database API authentication
* Ability to query backup/logs of addon
* Add missing Addon#Status field

## v1.5.2

* Remove os.Exit(), reliquat from split between CLI and client.
* Update wording
* Fix display of alert table

## v1.5.1

* Update deps

## v1.5.0

* Add AppsForceHTTPS
* Add AppsStickySession
* Add AppID in App subresources
* Collaborator.Status is now of type CollaboratorStatus, and constants are defined

## v1.4.1

* Add UserID to Collaborator

## v1.4.0

* Add Fullname for `User` model
* Ability to create an email notifier
* Access to one-off audit logs

## v1.3.2

* Add events NewAlert, Alert, DeleteAlert, NewAutoscaler, DeleteAutoscaler

## v1.3.0

* Change keys endpoint to point to the authentication service instead of the main API
* Add `GithubLinkService` implementation

## v1.2.0

* Refactoring, use interface instead of private struct

## v1.1.0

* API Token methods + authentication

## v1.0.0

* Initial tag
