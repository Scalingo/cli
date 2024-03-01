# Changelog

## To be Released

## v2.4.0

* docs(errors): deprecate use of `errgo` in `ErrCtx`

## v2.3.0

* feat: add `Is` and `As` to match standard library
* fix: remove `NoteMask`. `Notef` should be use instead.

## v2.2.0

* feat: add UnwrapError to unwrap one error.

## v2.1.0

* feat: add New function to `ErrCtx`
* feat: IsRootCause and RootCause are taking in account `ErrCtx` underlying errors
* feat: RootCtxOrFallback retrieves the deepest context from wrapped errors.

## v2.0.0

* fix: privatify `ErrgoRoot`
* build(deps): bump github.com/stretchr/testify from 1.8.0 to 1.8.1

## v1.1.1

* chore(go): use go 1.17
* build(deps): bump github.com/stretchr/testify from 1.7.0 to 1.7.1

## v1.1.0

* Bump go version to 1.16

## v1.0.0

* Initial breakdown of go-utils into subpackages
