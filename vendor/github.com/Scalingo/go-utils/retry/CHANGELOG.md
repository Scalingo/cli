# Changelog

## To be Released

## v1.4.2

* fix(retry): do not call the method if the context is already done

## v1.4.1

* chore(deps): bump multiple dependencies

## v1.4.0

* feat(retry): add exponential backoff wait duration
* refactor(retry): use `synctest.Test` for unit tests

## v1.3.0

* feat(unwrap): errors are now unwrappable by errors.Is/errors.As

## v1.2.2

* feat(logging): add `WithLoggingOnAttemptError` option allowing to log errors before the last attempt

## v1.2.1

* chore(go): corrective bump - Go version regression from 1.24.3 to 1.24

## v1.2.0

* chore(go): upgrade to Go 1.24
* build(deps): bump github.com/stretchr/testify from 1.8.0 to 1.8.1

## v1.1.1

* chore(go): use go 1.17
* build(deps): bump github.com/stretchr/testify from 1.7.0 to 1.7.1

## v1.1.0

* Bump go version to 1.16

## v1.0.0

* Initial breakdown of go-utils into subpackages
