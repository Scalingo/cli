package apps

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/pkg/errors"
)

func ForceHTTPS(appName string, enable bool) error {
	c := config.ScalingoClient()
	_, err := c.AppsForceHTTPS(appName, enable)
	if err != nil {
		return errors.Wrap(err, "fail to configure force-https feature")
	}

	var action string
	if enable {
		action = "enable"
	} else {
		action = "disable"
	}

	io.Statusf("Force HTTPS has been %sd on %s", action, appName)
	return nil
}
