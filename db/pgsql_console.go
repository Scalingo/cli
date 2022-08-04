package db

import (
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/config"
)

type PgSQLConsoleOpts struct {
	App  string
	Size string
}

func PgSQLConsole(opts PgSQLConsoleOpts) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	postgreSQLURL, user, _, err := dbURL(c, opts.App, "SCALINGO_POSTGRESQL", []string{"postgres://", "postgis://"})
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
