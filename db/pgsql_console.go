package db

import (
	"net"

	"github.com/Scalingo/cli/apps"
	"gopkg.in/errgo.v1" // "mysql2://" for ruby driver 'mysql2'
)

type PgSQLConsoleOpts struct {
	App  string
	Size string
}

func PgSQLConsole(opts PgSQLConsoleOpts) error {
	postgreSQLURL, user, password, err := dbURL(opts.App, "SCALINGO_POSTGRESQL", []string{"postgres://", "postgis://"})
	if err != nil {
		return errgo.Mask(err)
	}

	host, port, err := net.SplitHostPort(postgreSQLURL.Host)
	if err != nil {
		return errgo.Newf("%v has an invalid host", postgreSQLURL)
	}

	runOpts := apps.RunOpts{
		DisplayCmd: "pgsql-console " + user,
		App:        opts.App,
		Cmd:        []string{"dbclient-fetcher", "pgsql", "&&", "psql"},
		Size:       opts.Size,
		CmdEnv: []string{
			"PGHOST=" + host,
			"PGPORT=" + port,
			"PGUSER=" + user,
			"PGPASSWORD=" + password,
			"PGDATABASE=" + user,
		},
	}

	err = apps.Run(runOpts)
	if err != nil {
		return errgo.Newf("Fail to run PostgreSQL console: %v", err)
	}

	return nil
}
