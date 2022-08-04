package db

import (
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/apps"
	"github.com/Scalingo/cli/config"
)

type MongoConsoleOpts struct {
	App  string
	Size string
}

func MongoConsole(opts MongoConsoleOpts) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	mongoURL, _, _, err := dbURL(c, opts.App, "SCALINGO_MONGO", []string{"mongodb://"})
	if err != nil {
		return errgo.Mask(err)
	}

	command := []string{"dbclient-fetcher", "mongo", "&&", "mongo"}
	if mongoURL.Query().Get("ssl") == "true" {
		command = append(command, "--ssl", "--sslAllowInvalidCertificates")
	}

	command = append(command, "'"+mongoURL.String()+"'")

	runOpts := apps.RunOpts{
		DisplayCmd: "mongo-console",
		App:        opts.App,
		Cmd:        command,
		Size:       opts.Size,
	}

	err = apps.Run(runOpts)
	if err != nil {
		return errgo.Newf("Fail to run MongoDB console: %v", err)
	}

	return nil
}
