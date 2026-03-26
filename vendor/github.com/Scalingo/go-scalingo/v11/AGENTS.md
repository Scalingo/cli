# AGENTS Guidance

This repository enforces strict linting. Follow these rules for new and modified Go code.

## Imports And Formatting
- Always run `goimports -w -local github.com/Scalingo` on modified Go files.
- Keep imports in the order enforced by `goimports`:
  1. Standard library
  2. Third-party
  3. Local imports grouped under `github.com/Scalingo`

## Go Modernization
- Use `any` instead of `interface{}`.

## Tests
- Prefer `t.Context()` over `context.Background()` inside tests.
- In `http.HandlerFunc` test callbacks, do not use `require.*`.
- In handlers, use `assert.*` and return early when an error should stop assertions in that callback.
