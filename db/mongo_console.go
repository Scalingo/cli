package db

import (
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/apps"
)

type MongoConsoleOpts struct {
	App  string
	Size string
}

func MongoConsole(opts MongoConsoleOpts) error {
	mongoURL, _, _, err := dbURL(opts.App, "SCALINGO_MONGO", []string{"mongodb"})
	if err != nil {
		return errgo.Mask(err)
	}

	command := []string{"dbclient-fetcher", "mongo", "&&", "mongo"}
	if mongoURL.Query().Get("ssl") == "true" {
		command = append(command, "--ssl", "--sslAllowInvalidCertificates")
	}

	err = apps.Run(apps.RunOpts{
		DisplayCmd: "mongo-console",
		App:        opts.App,
		Cmd:        append(command, "'"+mongoURL.String()+"'"),
		Size:       opts.Size,
	})
	if err != nil {
		return errgo.Newf("Fail to run MongoDB console: %v", err)
	}

	return nil
}
