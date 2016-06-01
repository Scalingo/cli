package db

import (
	"net"

	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1" // "mysql2://" for ruby driver 'mysql2'
	"github.com/Scalingo/cli/apps"
)

func PgSQLConsole(app string) error {

	postgreSQLURL, user, password, err := dbURL(app, "SCALINGO_POSTGRESQL", []string{"postgres://", "postgis://"})
	if err != nil {
		return errgo.Mask(err)
	}

	host, port, err := net.SplitHostPort(postgreSQLURL.Host)
	if err != nil {
		return errgo.Newf("%v has an invalid host", postgreSQLURL)
	}

	opts := apps.RunOpts{
		DisplayCmd: "pgsql-console " + user,
		App:        app,
		Cmd:        []string{"psql"},
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
		return errgo.Newf("Fail to run PostgreSQL console: %v", err)
	}

	return nil
}
