[ ![Codeship Status for Scalingo/go-scalingo](https://app.codeship.com/projects/cf518dc0-0034-0136-d6b3-5a0245e77f67/status?branch=master)](https://app.codeship.com/projects/279805)

# Go client for Scalingo API v8.3.0

This repository is the Go client for the [Scalingo APIs](https://developers.scalingo.com/).

## Getting Started

```go
package main

import (
	"github.com/Scalingo/go-scalingo/v8"
)

func getClient() (*scalingo.Client, error) {
	config := scalingo.ClientConfig{
		APIEndpoint: "https://api.osc-fr1.scalingo.com", // Possible endpoints can be found at https://developers.scalingo.com/#endpoints
		APIToken: "tk-us-XYZXYZYZ", // You can create a token in the dashboard at Profile > Token > Create new token
	}
	return scalingo.New(config)
}

func main() {
	client, err := getClient()
	if err != nil {
		panic(err)
	}
	apps, err := client.AppsList()
	if err != nil {
		panic(err)
	}
	for _, app := range apps {
		println("App: " + app.Name)
	}
}
```

## Explore

As this Go client maps all the public Scalingo APIs, you can explore the [API documentation](https://developers.scalingo.com/).
If you need an implementation example, you can take a look at the [Scalingo CLI code](https://github.com/Scalingo/cli).

## Repository Processes

### Add Support for a New Event

A couple of files must be updated when adding support for a new event type. For
instance if your event type is named `my_event`:
* `events_struct.go`:
    * Add the `EventMyEvent` constant
    * Add the `EventMyEventTypeData` structure
    * Add `EventMyEventType` structure which embeds a field `TypeData` of the
        type `EventMyEventTypeData`.
    * Implement function `String` for `EventMyEventType`
    * [optional] Implement `Who` function for `EventMyEventType`. E.g. if the
        event type can be created by an addon.

Once the Event has been added run the following command to update boilerplate code

```
go generate
```

### Client HTTP Errors

HTTP errors are managed in the file
[http/errors.go](https://github.com/Scalingo/go-scalingo/blob/master/http/errors.go).
It follows the Scalingo standards detailed in the [developers
documentation](https://developers.scalingo.com/index#errors).

### Release a New Version

#### Create the new release

Bump new version number in:

- `CHANGELOG.md`
- `README.md`
- `version.go`

Commit, tag and create a new release:

```sh
version="8.3.0"

git switch --create release/${version}
git add CHANGELOG.md README.md version.go
git commit -m "Bump v${version}"
git push --set-upstream origin release/${version}
gh pr create --reviewer=EtienneM --title "$(git log -1 --pretty=%B)"
```

Once the pull request merged, you can tag the new release.

```sh
git tag v${version}
git push origin master v${version}
gh release create v${version}
```

The title of the release should be the version number and the text of the
release is the same as the changelog.

#### Document the new release on the documentation

When a new version is released, it must be documented on the general [changelog](https://doc.scalingo.com/changelog).
In order to do that, a new file must be created in this [folder](https://github.com/Scalingo/documentation/tree/master/src/changelog/sdk/_posts).
