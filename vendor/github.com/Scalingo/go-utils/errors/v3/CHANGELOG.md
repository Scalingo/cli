# Changelog

## To be Released

## v3.2.0

* feat(errors) Add `Join` wrapping `errors.Join` from standard library
* feat(errctx) `RootCtxOrFallback` is now compatible with `Join(...error)` and returns context of the first error

## v.3.1.1

* feat(errctx) RootCtxOrFallback is now compatible with different versions of errors package.

## v.3.1.0

* feat(UnwrapError) `UnwrapError` now unwraps errors which implement an `Unwrap()` method.

## v3.0.0

* fix(errors): `Build` returns `error` not `*ValidationErrors` (BREAKING CHANGE)
* refactor(errors): remove deprecated functions (BREAKING CHANGE)

## v2.5.1

* chore(go): corrective bump - Go version regression from 1.24.3 to 1.24

## v2.5.0

* chore(go): upgrade to Go 1.24

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
