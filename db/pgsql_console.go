package db

import (
	"context"

	"github.com/Scalingo/go-utils/errors/v3"

	"github.com/Scalingo/cli/apps"
)

type PgSQLConsoleOpts struct {
	App          string
	Size         string
	VariableName string
}

func PgSQLConsole(ctx context.Context, opts PgSQLConsoleOpts) error {
	if opts.VariableName == "" {
		opts.VariableName = "SCALINGO_POSTGRESQL"
	}
	postgreSQLURL, user, _, err := dbURL(ctx, opts.App, opts.VariableName, []string{"postgres", "postgis", "postgresql"})
	if err != nil {
		return errors.Wrapf(ctx, err, "resolve PostgreSQL URL from %s", opts.VariableName)
	}

	runOpts := apps.RunOpts{
		DisplayCmd: "pgsql-console " + user,
		App:        opts.App,
		Cmd:        []string{"dbclient-fetcher", "pgsql", "&&", "psql", "'" + postgreSQLURL.String() + "'"},
		Size:       opts.Size,
	}

	err = apps.Run(ctx, runOpts)
	if err != nil {
		return errors.Newf(ctx, "fail to run PostgreSQL console: %v", err)
	}

	return nil
}
