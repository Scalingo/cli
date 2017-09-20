package db

import (
	"github.com/Scalingo/cli/apps"
	"gopkg.in/errgo.v1" // "mysql2://" for ruby driver 'mysql2'
)

type PgSQLConsoleOpts struct {
	App  string
	Size string
}

func PgSQLConsole(opts PgSQLConsoleOpts) error {
	postgreSQLURL, user, _, err := dbURL(opts.App, "SCALINGO_POSTGRESQL", []string{"postgres://", "postgis://"})
	if err != nil {
		return errgo.Mask(err)
	}

	runOpts := apps.RunOpts{
		DisplayCmd: "pgsql-console " + user,
		App:        opts.App,
		Cmd:        []string{"dbclient-fetcher", "pgsql", "&&", "psql", "'" + postgreSQLURL.String() + "'"},
		Size:       opts.Size,
	}

	err = apps.Run(runOpts)
	if err != nil {
		return errgo.Newf("Fail to run PostgreSQL console: %v", err)
	}

	return nil
}
