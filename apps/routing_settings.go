package apps

import (
	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"gopkg.in/errgo.v1"
)

func ForceHTTPS(appName string, enable bool) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	_, err = c.AppsForceHTTPS(appName, enable)
	if err != nil {
		return errgo.Notef(err, "fail to configure force-https feature")
	}

	var action string
	if enable {
		action = "enable"
	} else {
		action = "disable"
	}

	io.Statusf("Force HTTPS has been %sd on %s\n", action, appName)
	return nil
}

func StickySession(appName string, enable bool) error {
	c, err := config.ScalingoClient()
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}
	_, err = c.AppsStickySession(appName, enable)
	if err != nil {
		return errgo.Notef(err, "fail to configure sticky-session feature")
	}

	var action string
	if enable {
		action = "enable"
	} else {
		action = "disable"
	}

	io.Statusf("Sticky session has been %sd on %s\n", action, appName)
	return nil
}
