package db

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/go-utils/errors/v3"
)

type RedisConsoleOpts struct {
	App          string
	Size         string
	VariableName string
}

func RedisConsole(ctx context.Context, opts RedisConsoleOpts) error {
	if opts.VariableName == "" {
		opts.VariableName = "SCALINGO_REDIS"
	}
	redisURL, _, password, err := dbURL(ctx, opts.App, opts.VariableName, []string{"redis", "rediss"})
	if err != nil {
		return errors.Wrapf(ctx, err, "resolve Redis URL from %s", opts.VariableName)
	}

	if redisURL.Scheme == "rediss" {
		return errors.New(ctx, "Redis console is not available when TLS connections are enforced")
	}

	host, port, err := net.SplitHostPort(redisURL.Host)
	if err != nil {
		return fmt.Errorf("%v has an invalid host", redisURL)
	}

	runOpts := apps.RunOpts{
		DisplayCmd:    "redis-console " + strings.Split(host, ".")[0],
		App:           opts.App,
		Cmd:           []string{"dbclient-fetcher", "redis", "&&", "redis-cli", "-h", host, "-p", port, "-a", password},
		Size:          opts.Size,
		StdinCopyFunc: valkeyStdinCopy,
	}

	err = apps.Run(ctx, runOpts)
	if err != nil {
		return fmt.Errorf("fail to run Redis console: %v", err)
	}

	return nil
}
