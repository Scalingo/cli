package db

import (
	"net"

	"github.com/Scalingo/cli/apps"
	"gopkg.in/errgo.v1"
)

func PgSQLConsole(app string) error {
	// "mysql2://" for ruby driver 'mysql2'
	postgreSQLURL, user, password, err := dbURL(app, "POSTGRESQL", []string{"postgres://"})
	if err != nil {
		return errgo.Mask(err)
	}

	host, port, err := net.SplitHostPort(postgreSQLURL.Host)
	if err != nil {
		return errgo.Newf("%v has an invalid host", postgreSQLURL)
	}

	opts := apps.RunOpts{
		App: app,
		Cmd: []string{"psql"},
		CmdEnv: []string{
			"PGHOST=" + host,
			"PGPORT=" + port,
			"PGUSER=" + user,
			"PGPASSWORD=" + password,
			"PGDATABASE=" + user,
		},
	}

	err = apps.Run(opts)
	if err != nil {
		return errgo.Newf("Fail to run redis console: %v", err)
	}

	return nil
}
