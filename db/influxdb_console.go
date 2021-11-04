package db

import (
	"net"
	"strings"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/config"
)

type InfluxDBConsoleOpts struct {
	App  string
	Size string
}

func InfluxDBConsole(opts InfluxDBConsoleOpts) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	influxdbURL, username, password, err := dbURL(c, opts.App, "SCALINGO_INFLUX", []string{"http://", "https://"})
	if err != nil {
		return errgo.Mask(err)
	}

	host, port, err := net.SplitHostPort(influxdbURL.Host)
	if err != nil {
		return errgo.Newf("%v has an invalid host", influxdbURL)
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

	err = apps.Run(runOpts)
	if err != nil {
		return errgo.Newf("Fail to run InfluxDB console: %v", err)
	}

	return nil
}
