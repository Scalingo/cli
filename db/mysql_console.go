package db

import (
	"fmt"
	"net"

	"github.com/Scalingo/cli/apps"
	"gopkg.in/errgo.v1"
)

func MySQLConsole(app string) error {
	// "mysql2://" for ruby driver 'mysql2'
	mySQLURL, user, password, err := dbURL(app, "MYSQL", []string{"mysql://", "mysql2://"})
	if err != nil {
		return errgo.Mask(err)
	}

	host, port, err := net.SplitHostPort(mySQLURL.Host)
	if err != nil {
		return errgo.Newf("%v has an invalid host", mySQLURL)
	}

	opts := apps.RunOpts{
		App: app,
		Cmd: []string{"mysql", "-h", host, "-P", port, fmt.Sprintf("--password=%v", password), "-u", user, user},
	}

	err = apps.Run(opts)
	if err != nil {
		return errgo.Newf("Fail to run redis console: %v", err)
	}

	return nil
}
