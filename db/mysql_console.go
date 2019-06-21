package db

import (
	"fmt"
	"net"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/config"
	"gopkg.in/errgo.v1" // "mysql2://" for ruby driver 'mysql2'
)

type MySQLConsoleOpts struct {
	App  string
	Size string
}

func MySQLConsole(opts MySQLConsoleOpts) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	mySQLURL, user, password, err := dbURL(c, opts.App, "SCALINGO_MYSQL", []string{"mysql://", "mysql2://"})
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
		return errgo.Newf("Fail to run MySQL console: %v", err)
	}

	return nil
}
