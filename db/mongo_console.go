package db

import (
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

	runOpts := apps.RunOpts{
		DisplayCmd: "mongo-console " + user,
		App:        opts.App,
		Cmd:        []string{"mongo", "-u", user, "-p", password, mongoURL.Host + "/" + user},
		Size:       opts.Size,
	}

	err = apps.Run(runOpts)
	if err != nil {
		return errgo.Newf("Fail to run MongoDB console: %v", err)
	}

	return nil
}
