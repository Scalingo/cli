# Package `retry` v1.1.1

This library implements a retryer: a generic way to execute some code at
regular interval.

## Usage

Creating a new retryer is done with:

```go
retry.New(retry.WithWaitDuration(10*time.Second), retry.WithMaxAttempts(5))
```

Possible options for the constructor are:

- `WithWaitDuration`: time interval between each execution of the code
  (default to 10 seconds).
- `WithMaxDuration`: the retryer will stop executing after the specified
  amount of time (disabled by default).
- `WithMaxAttempts`: the retryer will execute the at most N times. N being
  this max attempts parameter (default to 5).
- `WithoutMaxAttempts`: disable the max attempts parameter.

Then execute the retryer with:

```go
err = retryer.Do(ctx, func(ctx context.Context) error {
    // ...
	})
```

The function given as parameter only returns an error. If this error has the
type `RetryCancelError`, it stops the execution of the retryer. You can
create such error with `NewRetryCancelError`.
