package db

import (
	"net"
	"strings"

	"github.com/Scalingo/cli/apps"
	"gopkg.in/errgo.v1"
)

type InfluxDBConsoleOpts struct {
	App  string
	Size string
}

func InfluxDBConsole(opts InfluxDBConsoleOpts) error {
	influxdbURL, username, password, err := dbURL(opts.App, "SCALINGO_INFLUX", []string{"http://"})
	if err != nil {
		return errgo.Mask(err)
	}

	host, port, err := net.SplitHostPort(influxdbURL.Host)
	if err != nil {
		return errgo.Newf("%v has an invalid host", influxdbURL)
	}

	runOpts := apps.RunOpts{
		DisplayCmd: "influxdb-console " + strings.Split(host, ".")[0],
		App:        opts.App,
		Cmd:        []string{"influx", "-host", host, "-port", port, "-username", username, "-password", password},
		Size:       opts.Size,
	}

	err = apps.Run(runOpts)
	if err != nil {
		return errgo.Newf("Fail to run InfluxDB console: %v", err)
	}

	return nil
}
