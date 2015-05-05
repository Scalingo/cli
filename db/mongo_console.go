package db

import (
	"github.com/Scalingo/cli/apps"
	"gopkg.in/errgo.v1"
)

func MongoConsole(app string) error {
	// "mysql2://" for ruby driver 'mysql2'
	mongoURL, user, password, err := dbURL(app, "MONGO", []string{"mongodb://"})
	if err != nil {
		return errgo.Mask(err)
	}

	opts := apps.RunOpts{
		App: app,
		Cmd: []string{"mongo", "-u", user, "-p", password, mongoURL.Host + "/" + user},
	}

	err = apps.Run(opts)
	if err != nil {
		return errgo.Newf("Fail to run redis console: %v", err)
	}

	return nil
}
