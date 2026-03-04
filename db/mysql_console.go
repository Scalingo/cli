package db

import (
	"context"
	"fmt"
	"net"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/go-utils/errors/v3"
)

type MySQLConsoleOpts struct {
	App          string
	Size         string
	VariableName string
}

func MySQLConsole(ctx context.Context, opts MySQLConsoleOpts) error {
	if opts.VariableName == "" {
		opts.VariableName = "SCALINGO_MYSQL"
	}
	mySQLURL, user, password, err := dbURL(ctx, opts.App, opts.VariableName, []string{"mysql", "mysql2"})
	if err != nil {
		return errors.Wrapf(ctx, err, "resolve MySQL URL from %s", opts.VariableName)
	}

	host, port, err := net.SplitHostPort(mySQLURL.Host)
	if err != nil {
		return errors.Newf(ctx, "%v has an invalid host", mySQLURL)
	}

	runOpts := apps.RunOpts{
		DisplayCmd: "mysql-console " + user,
		App:        opts.App,
		Cmd:        []string{"dbclient-fetcher", "mysql", "&&", "mysql", "-h", host, "-P", port, fmt.Sprintf("--password=%v", password), "-u", user, user},
		Size:       opts.Size,
	}

	err = apps.Run(ctx, runOpts)
	if err != nil {
		return errors.Newf(ctx, "fail to run MySQL console: %v", err)
	}

	return nil
}
