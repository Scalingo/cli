package db

import (
	"context"
	"net"
	"strings"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/go-utils/errors/v3"
)

type InfluxDBConsoleOpts struct {
	App          string
	Size         string
	VariableName string
}

func InfluxDBConsole(ctx context.Context, opts InfluxDBConsoleOpts) error {
	if opts.VariableName == "" {
		opts.VariableName = "SCALINGO_INFLUX"
	}
	influxdbURL, username, password, err := dbURL(ctx, opts.App, opts.VariableName, []string{"http", "https"})
	if err != nil {
		return errors.Wrapf(ctx, err, "resolve InfluxDB URL from %s", opts.VariableName)
	}

	host, port, err := net.SplitHostPort(influxdbURL.Host)
	if err != nil {
		return errors.Newf(ctx, "%v has an invalid host", influxdbURL)
	}

	cmd := []string{"dbclient-fetcher", "influxdb", "&&", "influx"}

	if influxdbURL.Scheme == "https" {
		cmd = append(cmd, "-ssl", "-unsafeSsl")
	}

	cmd = append(cmd, "-host", host, "-port", port, "-username", username, "-password", password, "-database", influxdbURL.Path[1:])

	runOpts := apps.RunOpts{
		DisplayCmd: "influxdb-console " + strings.Split(host, ".")[0],
		App:        opts.App,
		Cmd:        cmd,
		Size:       opts.Size,
	}

	err = apps.Run(ctx, runOpts)
	if err != nil {
		return errors.Newf(ctx, "fail to run InfluxDB console: %v", err)
	}

	return nil
}
