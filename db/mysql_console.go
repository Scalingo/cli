package db

import (
	"fmt"
	"net"

	"gopkg.in/errgo.v1" // "mysql2://" for ruby driver 'mysql2'
	"github.com/Scalingo/cli/apps"
)

func MySQLConsole(app string) error {

	mySQLURL, user, password, err := dbURL(app, "SCALINGO_MYSQL", []string{"mysql://", "mysql2://"})
	if err != nil {
		return errgo.Mask(err)
	}

	host, port, err := net.SplitHostPort(mySQLURL.Host)
	if err != nil {
		return errgo.Newf("%v has an invalid host", mySQLURL)
	}

	opts := apps.RunOpts{
		DisplayCmd: "mysql-console " + user,
		App:        app,
		Cmd:        []string{"mysql", "-h", host, "-P", port, fmt.Sprintf("--password=%v", password), "-u", user, user},
	}

	err = apps.Run(opts)
	if err != nil {
		return errgo.Newf("Fail to run MySQL console: %v", err)
	}

	return nil
}
