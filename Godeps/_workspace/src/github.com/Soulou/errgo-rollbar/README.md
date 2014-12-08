## Rollbar - Errgo binding

Rollbar Go client: `https://github.com/stvp/rollbar`
Errgo: `https://github.com/go-errgo/errgo`

### Errgo masking system

Errgo is based on a error masking system which gather information
each time an error is masked and returned to the called.

From this stack of masked error, we can build the calling path of
the error and get the precise origin of a given error.

### Example

```go
func A() error {
  return errgo.Mask(err)
}

func B() error {
  return errgo.Mask(A())
}

func main() {
  err := B()
  rollbar.ErrorWithStack(rollbar.ERR, err, errgorollbar.BuildStack(err))
  rollbar.Wait()
}
```

Then your are able to see on your rollbal dashboard that the correct
stack trace has been sent.
