# Changelog

## To Be Released

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
