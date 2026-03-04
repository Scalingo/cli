package apps

import (
	"context"
	"fmt"

	"github.com/pkg/browser"

	"github.com/Scalingo/go-utils/errors/v3"
)

func Dashboard(appName string, region string) error {
	url := fmt.Sprintf("https://dashboard.scalingo.com/apps/%s/%s", region, appName)

	if err := browser.OpenURL(url); err != nil {
		return errors.Wrapf(context.Background(), err, "fail to open dashboard in browser")
	}

	return nil
}
