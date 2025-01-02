# Package `logger` v1.3.1

This package will provide you a generic way to handle logging.

## Configuration

This plugin will configure himself automatically using the following environment variables:

 * `LOGGER_TYPE`: define the logger output type (values: `json`, `text`) (default: `text`)
 * `LOGGER_LEVEL`: define the minimum output level of the logger (values: `panic`, `fatal`, `warn`, `info`, `debug`) (default: `info`)

## Usage

```go
log := logger.Default() // Return a default logger
```

## Context

The logger can be passed in a context so it can retain fields.

```go
func main() {
  log := logger.Default().WithFields(logrus.Fields{"caller": "main"})
  add(logger.ToCtx(context.Background(), log), 1, 2)
}

def add(ctx context.Context, a, b int) int {
  log := logger.Get(ctx)
  log.Info("Starting add operation")

  log.WithField("operation", "add")
  do(logger.ToCtx(ctx, log), a,b, func(a,b int)int{return a+b})
}

def do(ctx context.Context, a,b int, op fun(int, int)int) {
  log := logger.Get(ctx)
  log.Info("Doing operation")
  op(a,b)
}
```

```shell

2017-08-27 11:10:10 [INFO] Starting add operation caller=main
2017-08-27 11:10:10 [INFO] Do operation caller=main operation=add
```

## Plugins

This logger accept plugins which can register hooks on the logger.

### Known plugins

* [rollbar](https://github.com/Scalingo/go-utils/tree/master/logger/plugins/rollbarplugin)
