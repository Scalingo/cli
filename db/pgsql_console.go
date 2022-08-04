package db

import (
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/apps"
)

type PgSQLConsoleOpts struct {
	App          string
	Size         string
	VariableName string
}

func PgSQLConsole(opts PgSQLConsoleOpts) error {
	if opts.VariableName == "" {
		opts.VariableName = "SCALINGO_POSTGRESQL"
	}
	postgreSQLURL, user, _, err := dbURL(opts.App, opts.VariableName, []string{"postgres://", "postgis://"})
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
		return errgo.Newf("fail to run PostgreSQL console: %v", err)
	}

	return nil
}
