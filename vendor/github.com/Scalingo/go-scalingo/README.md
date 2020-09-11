[ ![Codeship Status for Scalingo/go-scalingo](https://app.codeship.com/projects/cf518dc0-0034-0136-d6b3-5a0245e77f67/status?branch=master)](https://app.codeship.com/projects/279805)

# Go client for Scalingo API v4.7.2

This repository is the Go client for the [Scalingo APIs](https://developers.scalingo.com/).

## Getting Started

```go
package main

import (
	"github.com/Scalingo/go-scalingo"
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
    * Add support for this event type in the `Specialize` function
    * [optional] Implement `Who` function for `EventMyEventType`. E.g. if the
        event type can be created by an addon.
* `events_boilerplate.go`: implement the `TypeDataPtr` function for the new
    `EventMyEventType` structure.

### Client HTTP Errors

HTTP errors are managed in the file
[http/errors.go](https://github.com/Scalingo/go-scalingo/blob/master/http/errors.go).
It follows the Scalingo standards detailed in the [developers
documentation](https://developers.scalingo.com/index#errors).

### Release a New Version

Bump new version number in:

- `CHANGELOG.md`
- `README.md`
- `version.go`

Commit, tag and create a new release:

```sh
git add CHANGELOG.md README.md version.go
git commit -m "Bump v4.7.1"
git tag v4.7.1
git push origin master
git push --tags
hub release create v4.7.1
```

The title of the release should be the version number and the text of the
release is the same as the changelog.
