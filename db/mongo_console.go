package db

import (
	"fmt"

	"github.com/Scalingo/cli/apps"
	"gopkg.in/errgo.v1" // "mysql2://" for ruby driver 'mysql2'
)

type MongoConsoleOpts struct {
	App  string
	Size string
}

func MongoConsole(opts MongoConsoleOpts) error {
	mongoURL, user, password, err := dbURL(opts.App, "SCALINGO_MONGO", []string{"mongodb://"})
	if err != nil {
		return errgo.Mask(err)
	}

	host := mongoURL.Host
	if mongoURL.Query().Get("replicaSet") != "" {
		host = fmt.Sprintf("%s/%s", mongoURL.Query().Get("replicaSet"), mongoURL.Host)
	}
	runOpts := apps.RunOpts{
		DisplayCmd: "mongo-console " + user,
		App:        opts.App,
		Cmd:        []string{"dbclient-fetcher", "mongo", "&&", "mongo", "-u", user, "-p", password, host + "/" + user},
		Size:       opts.Size,
	}

	err = apps.Run(runOpts)
	if err != nil {
		return errgo.Newf("Fail to run MongoDB console: %v", err)
	}

	return nil
}
