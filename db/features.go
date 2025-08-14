package db

import (
	"context"
	"time"

	"github.com/briandowns/spinner"
	"github.com/urfave/cli/v3"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v8"
)

// EnableFeature is the command handler to enable a database feature on a given
// database addon, like 'force-ssl' or 'public-availability'
func EnableFeature(ctx context.Context, c *cli.Command, app, addon, feature string) error {
	spinner := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	spinner.Suffix = " Enabling database feature"
	spinner.Start()
	defer spinner.Stop()

	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	res, err := client.DatabaseEnableFeature(ctx, app, addon, feature)
	if err != nil {
		return errgo.Notef(err, "fail to enable feature '%v'", feature)
	}
	spinner.Stop()

	switch res.Status {
	case scalingo.DatabaseFeatureStatusActivated:
		io.Statusf("Feature %v has been enabled\n", feature)
	case scalingo.DatabaseFeatureStatusFailed:
		io.Warningf("Feature %v failed to get activated, please contact our support\n", feature)
	case scalingo.DatabaseFeatureStatusPending:
		io.Statusf("Feature %v is being enabled\n", feature)
	}

	if res.Status == scalingo.DatabaseFeatureStatusPending && c.Bool("synchronous") {
		io.Infof("Waiting for operation completion...")
		err = waitFeatureUntilActivated(ctx, client, app, addon, feature)
		if err != nil {
			return errgo.Notef(err, "fail to wait for feature '%v' to be enabled", feature)
		}
	}

	return nil
}

func waitFeatureUntilActivated(ctx context.Context, client *scalingo.Client, app, addon, feature string) error {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		db, err := client.DatabaseShow(ctx, app, addon)
		if err != nil {
			return errgo.Notef(err, "fail to refresh database metadata")
		}
		for _, f := range db.Features {
			if f.Name == feature && f.Status != scalingo.DatabaseFeatureStatusPending {
				switch f.Status {
				case scalingo.DatabaseFeatureStatusActivated:
					io.Statusf("Feature %v has been activated\n", feature)
					return nil
				case scalingo.DatabaseFeatureStatusFailed:
					io.Warningf("Feature %v failed to get activated, please contact our support\n", feature)
					return nil
				}
			}
		}
	}
	return nil
}

// DisableFeature is the command handler to disable a database feature on a
// database addon like 'force-ssl' or 'public-availability'
func DisableFeature(ctx context.Context, app, addon, feature string) error {
	spinner := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	spinner.Suffix = " Disabling database feature"
	spinner.Start()
	defer spinner.Stop()

	client, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	_, err = client.DatabaseDisableFeature(ctx, app, addon, feature)
	if err != nil {
		return errgo.Notef(err, "fail to disable feature '%v'", feature)
	}
	spinner.Stop()

	io.Statusf("Feature %v has been disabled.\n", feature)

	return nil
}
