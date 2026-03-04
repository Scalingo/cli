package db

import (
	"context"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/go-utils/errors/v3"
)

type MongoConsoleOpts struct {
	App          string
	Size         string
	VariableName string
}

func MongoConsole(ctx context.Context, opts MongoConsoleOpts) error {
	if opts.VariableName == "" {
		opts.VariableName = "SCALINGO_MONGO"
	}
	mongoURL, _, _, err := dbURL(ctx, opts.App, opts.VariableName, []string{"mongodb"})
	if err != nil {
		return errors.Wrapf(ctx, err, "resolve MongoDB URL from %s", opts.VariableName)
	}

	command := []string{"dbclient-fetcher", "mongo", "&&", "mongo"}
	if mongoURL.Query().Get("ssl") == "true" {
		command = append(command, "--ssl", "--sslAllowInvalidCertificates")
	}

	err = apps.Run(ctx, apps.RunOpts{
		DisplayCmd: "mongo-console",
		App:        opts.App,
		Cmd:        append(command, "'"+mongoURL.String()+"'"),
		Size:       opts.Size,
	})
	if err != nil {
		return errors.Newf(ctx, "fail to run MongoDB console: %v", err)
	}

	return nil
}
