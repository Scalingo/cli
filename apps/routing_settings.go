package apps

import (
	"context"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-utils/errors/v3"
)

func ForceHTTPS(ctx context.Context, appName string, enable bool) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}

	_, err = c.AppsForceHTTPS(ctx, appName, enable)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to configure force-https feature")
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

func StickySession(ctx context.Context, appName string, enable bool) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}
	_, err = c.AppsStickySession(ctx, appName, enable)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to configure sticky-session feature")
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

func RouterLogs(ctx context.Context, appName string, enable bool) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}

	_, err = c.AppsRouterLogs(ctx, appName, enable)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to configure router-logs feature")
	}

	var action string
	if enable {
		action = "enable"
	} else {
		action = "disable"
	}

	io.Statusf("Router logs have been %sd on %s\n", action, appName)
	return nil
}
