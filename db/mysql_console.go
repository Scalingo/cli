package db

import (
	"fmt"
	"net"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/apps"
)

type MySQLConsoleOpts struct {
	App          string
	Size         string
	VariableName string
}

func MySQLConsole(opts MySQLConsoleOpts) error {
	if opts.VariableName == "" {
		opts.VariableName = "SCALINGO_MYSQL"
	}
	mySQLURL, user, password, err := dbURL(opts.App, opts.VariableName, []string{"mysql", "mysql2"})
	if err != nil {
		return errgo.Mask(err)
	}

	host, port, err := net.SplitHostPort(mySQLURL.Host)
	if err != nil {
		return errgo.Newf("%v has an invalid host", mySQLURL)
	}

	runOpts := apps.RunOpts{
		DisplayCmd: "mysql-console " + user,
		App:        opts.App,
		Cmd:        []string{"dbclient-fetcher", "mysql", "&&", "mysql", "-h", host, "-P", port, fmt.Sprintf("--password=%v", password), "-u", user, user},
		Size:       opts.Size,
	}

	err = apps.Run(runOpts)
	if err != nil {
		return errgo.Newf("fail to run MySQL console: %v", err)
	}

	return nil
}
